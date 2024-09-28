package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/h3th-IV/climateer/pkg/database"
	"github.com/h3th-IV/climateer/pkg/model"
	"github.com/h3th-IV/climateer/pkg/utils"
	"go.uber.org/zap"
)

var _ http.Handler = &registerHandler{}

type registerHandler struct {
	logger      *zap.Logger
	mysqlClient database.Database
}

func NewRegistrationHandler(logger *zap.Logger, mysqlclient database.Database) *registerHandler {
	return &registerHandler{
		logger:      logger,
		mysqlClient: mysqlclient,
	}
}

func (handler *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		RegisterUser *model.Registration
	)

	if err := json.NewDecoder(r.Body).Decode(&RegisterUser); err != nil {
		resp["err"] = "unable to process request"
		handler.logger.Error("err decoding JSON object", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusNotFound)
		return
	}
	// there should be a frontend validation for all fields
	// the backend would assist to catch empty fields if the
	// frontend validation is compromised.
	if RegisterUser.FirstName == "" || RegisterUser.LastName == "" || RegisterUser.Email == "" || RegisterUser.Password == "" || RegisterUser.Phone == "" || RegisterUser.EduInstitute == "" {
		handler.logger.Error("some fields are empty")
		resp["err"] = "some fields are empty"
		apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusBadRequest)
		return
	}

	specialchars := strings.ContainsAny(RegisterUser.Password, "$ % @ !")
	passwdcount := len(RegisterUser.Password)
	if !specialchars {
		handler.logger.Error("password must contain special characters")
		resp["err"] = "password must contain special characters"
		apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusBadRequest)
		return
	}
	if passwdcount <= 7 {
		handler.logger.Error("password must contain at least 8 characters")
		resp["err"] = "password must contain at least 8 characters"
		apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusBadRequest)
		return
	}

	hashed_password, err := utils.HashPassword(RegisterUser.Password)
	if err != nil {
		handler.logger.Error("cannot hash password", zap.String("hashed password error", err.Error()))
		fmt.Printf("cannot hash password(%s)", RegisterUser.Password)
		return
	}
	newSessionKey := createSessionKey(RegisterUser.Email, time.Now())
	sanitize_email, err := utils.ValidateEmail(RegisterUser.Email)
	if err != nil || !sanitize_email {
		handler.logger.Error("email was malformed!", zap.Error(err))
		return
	}

	createUser, err := handler.mysqlClient.CreateUser(r.Context(), RegisterUser.FirstName, RegisterUser.LastName, RegisterUser.Email, hashed_password, RegisterUser.Phone, RegisterUser.EduInstitute, newSessionKey)
	if err != nil || !createUser {
		resp["err"] = "cannot register user, try again"
		handler.logger.Error("could not create user", zap.Any("error", err))
		apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusInternalServerError)
		return
	}
	resp["first_name"] = RegisterUser.FirstName
	resp["last_name"] = RegisterUser.LastName
	resp["email"] = RegisterUser.Email
	resp["phone"] = RegisterUser.Phone
	resp["edu_institute"] = RegisterUser.EduInstitute
	resp["session_key"] = newSessionKey
	apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusOK)
}
