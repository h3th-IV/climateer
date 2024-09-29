package handlers

import (
	"net/http"

	"github.com/h3th-IV/climateer/pkg/database"
	"go.uber.org/zap"
)

var _ http.Handler = &homeHandler{}

type homeHandler struct {
	logger      *zap.Logger
	mysqlClient database.Database
}

func NewHomeHandler(logger *zap.Logger, mysqlClient database.Database) *homeHandler {
	return &homeHandler{
		logger:      logger,
		mysqlClient: mysqlClient,
	}
}

func (handler *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var TTL = 60
	var resp = map[string]interface{}{}
	resp["message"] = "Welcome to Lasu Climate Repo"
	handler.logger.Info("Home Page was hit successfully")
	apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusOK)
}
