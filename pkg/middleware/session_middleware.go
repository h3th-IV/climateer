package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/h3th-IV/climateer/pkg/database"
	"github.com/h3th-IV/climateer/pkg/model"
	"github.com/h3th-IV/climateer/pkg/utils"
	"go.uber.org/zap"
)

type SessionMiddleware struct {
	logger      *zap.Logger
	mysqlclient database.Database
}

func NewSessionMiddleware(logger *zap.Logger, mysqlclient database.Database) *SessionMiddleware {
	return &SessionMiddleware{
		logger:      logger,
		mysqlclient: mysqlclient,
	}
}

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Err     string      `json:"err,omitempty"`
	Success bool        `json:"success"`
	TTL     int         `json:"ttl"`
}

func GetSuccessResponse(data interface{}, ttl int) []byte {
	resp := &response{
		Success: true,
		TTL:     ttl,
	}

	if data == nil || (reflect.ValueOf(data).Kind() == reflect.Ptr && reflect.ValueOf(data).IsNil()) {
		resp.Data = nil
	} else {
		resp.Data = data
	}
	responseBytes, _ := json.Marshal(resp)
	return responseBytes
}

func (smw *SessionMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionkey = r.FormValue("session_key")
		var sessionresp = map[string]interface{}{}
		uInfo := &model.User{}
		getBySession, err := smw.mysqlclient.GetBySessionKey(r.Context(), sessionkey)
		if err != nil {
			sessionresp["err"] = "no session key passed"
			smw.logger.Debug("cannot fetch user by sessionkey", zap.Any("fetch user error", err))
			w.Write(GetSuccessResponse(sessionresp["err"], 30))
			return
		}
		uInfo.ID = getBySession.ID
		uInfo.FirstName = getBySession.FirstName
		uInfo.LastName = getBySession.LastName
		uInfo.Email = getBySession.Email
		uInfo.Password = getBySession.Password
		uInfo.Phone = getBySession.Phone
		uInfo.EduInstitute = getBySession.EduInstitute
		uInfo.SessionKey = getBySession.SessionKey
		ctx := model.NewContext(r.Context(), uInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func middlewareResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
func JWTAuthRoutes(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fetch token from request
		AuthToken := r.Header.Get("Authorization")
		if AuthToken == "" {
			middlewareResponse(w, "please login to access resources", http.StatusUnauthorized)
			utils.Logger.Error("no auth token was provided")
			return
		}
		jwToken := strings.Split(AuthToken, " ")[1]

		token, err := jwt.Parse(jwToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		//check err
		if err != nil {
			if strings.Contains(err.Error(), "Token is expired") {
				middlewareResponse(w, "session expired, please sign in", http.StatusUnauthorized)
				utils.Logger.Error("err parsing token or invalid token", zap.Error(err))
				return
			}
			middlewareResponse(w, "please signIn to access resources", http.StatusUnauthorized)
			utils.Logger.Error("err parsing token or invalid token", zap.Error(err))
			return
		}

		if !token.Valid {
			middlewareResponse(w, "please signIn to access resources", http.StatusUnauthorized)
			utils.Logger.Error("invalid token")
			return
		}
		//check claims
		tokenClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			middlewareResponse(w, "please signIn to access resources", http.StatusUnauthorized)
			utils.Logger.Error("invalid token claims", zap.Bool("", ok))
			return
		}

		sessionKey, ok := tokenClaims["user"]
		if !ok {
			middlewareResponse(w, "please signIn to access resources", http.StatusUnauthorized)
			utils.Logger.Error("user is not Authorized")
			return
		}

		//store token(user.sessionKey)
		ctx := context.WithValue(r.Context(), utils.UserIDKey, sessionKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// auth routes
func AuthRoute(next http.Handler) http.Handler {
	return JWTAuthRoutes(next, utils.MYSTIC)
}
