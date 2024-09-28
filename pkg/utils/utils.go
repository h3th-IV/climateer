package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/h3th-IV/climateer/pkg/database"
	"github.com/h3th-IV/climateer/pkg/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return password, err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateEmail(email string) (bool, error) {
	if len(email) > 50 {
		return false, fmt.Errorf("email exceeds required length")
	}
	parts := strings.Split(strings.ToLower(email), "@")
	if len(parts) != 2 {
		return false, fmt.Errorf("email must contain @ symbol")
	}
	local, domain := parts[0], parts[1]
	if len(local) == 0 ||
		len(domain) == 0 {
		return false, fmt.Errorf("local or domain cannot be empty")
	}
	prev_char := rune(0)
	for _, char := range local {
		if strings.ContainsRune("!#$%&'*+-/=?^_`{|}~.", char) {
			if char == prev_char && char != '-' {
				return false, fmt.Errorf("cannot contain special chars before domain")
			}
		}
		prev_char = char
	}
	if strings.ContainsAny(email, " ") {
		return false, fmt.Errorf("email cannot contain spaces")
	}
	if len(local) > 64 || len(domain) > 255 {
		return false, fmt.Errorf("local part or domain part length exceeds the limit in the email")
	}
	return true, nil
}

var Logger, _ = zap.NewDevelopment()

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Logger.Info((fmt.Sprintf("%v - %v %v %v", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())))
		next.ServeHTTP(w, r)
	})
}

type mapKey string

const (
	UserIDKey mapKey = "user_id"
)

func AuthenticateUser(ctx context.Context, logger *zap.Logger, mysqlclient database.Database) (*model.User, error) {
	sessionKey, ok := ctx.Value(UserIDKey).(string)
	if !ok || sessionKey == "" {
		logger.Error("session key is missing")
		return nil, errors.New("please sign in to access this page")
	}

	user, err := mysqlclient.GetBySessionKey(ctx, sessionKey)
	if err != nil {
		logger.Error("user is not authorized", zap.Error(err))
		return nil, errors.New("please sign in to access this page")
	}

	return user, nil
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//if panic close connection
				w.Header().Set("Connection", "Close")
				//write internal server error
				ServerError(w, "Connection Closed inabruptly", fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ServerError(w http.ResponseWriter, errMsg string, err error) {
	errTrace := fmt.Sprintf("%v\n%v", err.Error(), debug.Stack())
	Logger.Error(errTrace)
	http.Error(w, errMsg, http.StatusInternalServerError)
}

func GenerateToken(user *model.User, expiry time.Duration, issuer, secret string) (string, error) {
	//set token expiry time
	bestBefore := time.Now().Add(expiry)
	//set token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  user.Email,
		"user":   user.SessionKey,
		"exp":    bestBefore.Unix(),
		"issuer": issuer,
	})
	//generate jwt token str and sign with secret key
	JWToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return JWToken, nil
}

var (
	MYSTIC    string
	JWTISSUER string
)
