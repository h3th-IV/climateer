package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/h3th-IV/climateer/pkg/database"
	"github.com/h3th-IV/climateer/pkg/model"
	"github.com/h3th-IV/climateer/pkg/utils"
	"go.uber.org/zap"
)

var _ http.Handler = &loginHandler{}

type loginHandler struct {
	logger      *zap.Logger
	mysqlclient database.Database
}

func NewLoginHandler(logger *zap.Logger, mysqlclient database.Database) *loginHandler {
	return &loginHandler{
		logger:      logger,
		mysqlclient: mysqlclient,
	}
}
func (handler *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		login *model.LoginCredentials
		TTL   = 60
		resp  = map[string]interface{}{}
	)

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		resp["err"] = "unable to process request"
		handler.logger.Error("err decoding request body", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}

	// Validate email and password fields
	if login.Email == "" || login.Password == "" {
		resp["err"] = "Email and password are required"
		handler.logger.Warn("Email or password missing in request")
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusBadRequest)
		return
	}

	checkuser, err := handler.mysqlclient.GetUserByEmail(r.Context(), login.Email)
	if err != nil {
		resp["err"] = "user does not exist"
		handler.logger.Error("user does not exist", zap.Any("checkuser", err))
		apiResponse(w, GetErrorResponseBytes(resp["err"], TTL, nil), http.StatusNotFound)
		return
	}
	if checkuser != nil {
		if checkuser.ID > 0 {
			handler.logger.Debug("found user", zap.Bool("user found", true))
			_ = utils.CheckPasswordHash(login.Password, checkuser.Password)
			loginnow, err := handler.mysqlclient.CheckUser(r.Context(), checkuser.Email, checkuser.Password)
			if err != nil {
				resp["err"] = "email or password incorrect"
				handler.logger.Error("email or password incorrect", zap.Any("login response", "email or password incorrect"))
				apiResponse(w, GetErrorResponseBytes(resp["err"], TTL, nil), http.StatusUnauthorized)
				return
			}
			if loginnow != nil {
				jwt, err := utils.GenerateToken(loginnow, 2*time.Hour, utils.JWTISSUER, utils.MYSTIC)
				if err != nil {
					resp["err"] = "unable to authenticate user"
					handler.logger.Error("err generating auth token")
					apiResponse(w, GetErrorResponseBytes(resp["err"], TTL, nil), http.StatusInternalServerError)
					return
				}
				resp["first_name"] = loginnow.FirstName
				resp["last_name"] = loginnow.LastName
				resp["email"] = loginnow.Email
				resp["phone"] = loginnow.Phone
				resp["edu_institute"] = loginnow.EduInstitute
				resp["jwt_token"] = jwt
				resp["message"] = "login successfull"
				apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusOK)
				handler.logger.Info("login successfull")
			}
		}
	}

}
