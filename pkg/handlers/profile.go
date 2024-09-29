package handlers

import (
	"net/http"

	"github.com/h3th-IV/climateer/pkg/database"
	"github.com/h3th-IV/climateer/pkg/utils"
	"go.uber.org/zap"
)

var _ http.Handler = &profileHandler{}

type profileHandler struct {
	logger      *zap.Logger
	mysqlclient database.Database
}

func NewProfileHandler(logger *zap.Logger, mysqlclient database.Database) *profileHandler {
	return &profileHandler{
		logger:      logger,
		mysqlclient: mysqlclient,
	}
}

func (handler *profileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var TTL = 60
	var resp = map[string]interface{}{}
	userInfo, err := utils.AuthenticateUser(r.Context(), handler.logger, handler.mysqlclient)
	if err != nil {
		resp["err"] = "please sign in to access this page"
		handler.logger.Warn("unauthorized user")
		apiResponse(w, GetErrorResponseBytes(resp["err"], TTL, nil), http.StatusUnauthorized)
		return
	}
	resp["id"] = userInfo.ID
	resp["first_name"] = userInfo.FirstName
	resp["last_name"] = userInfo.LastName
	resp["email"] = userInfo.Email
	resp["phone"] = userInfo.Phone
	resp["edu_institute"] = userInfo.EduInstitute
	apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusOK)
}
