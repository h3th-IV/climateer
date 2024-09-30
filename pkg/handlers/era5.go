package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/h3th-IV/climateer/pkg/database"
	"github.com/h3th-IV/climateer/pkg/model"
	"github.com/h3th-IV/climateer/pkg/utils"
	"go.uber.org/zap"
)

var _ http.Handler = &era5Handler{}

type era5Handler struct {
	logger      *zap.Logger
	mysqlclient database.Database
}

func NewEra5Handler(logger *zap.Logger, mysqlClient database.Database) *era5Handler {
	return &era5Handler{
		logger:      logger,
		mysqlclient: mysqlClient,
	}
}

func (handler *era5Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		TTL           = 60
		resp          = map[string]interface{}{}
		era5_request  model.ReanalysisRequest
		era5_response model.ReanalysisResponse
	)

	_, err := utils.AuthenticateUser(r.Context(), handler.logger, handler.mysqlclient)
	if err != nil {
		resp["err"] = "please sign in to access this page"
		handler.logger.Warn("unauthorized user")
		apiResponse(w, GetErrorResponseBytes(resp["err"], TTL, nil), http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&era5_request); err != nil {
		resp["err"] = "unable to process request"
		handler.logger.Error("err decoding json object", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusBadRequest)
		return
	}
	cmd := exec.Command("python3", "reanalysis.py", era5_request.Variable, era5_request.Year, era5_request.Month, era5_request.Day, era5_request.Time, era5_request.Area)
	output, err := cmd.Output()
	if err != nil {
		resp["err"] = "unable to proceed"
		handler.logger.Error("err executing api request", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(output, &era5_response); err != nil {
		resp["err"] = "unable to process response"
		handler.logger.Error("err parsing api response", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}

	publicDir := "./data"
	err = os.MkdirAll(publicDir, os.ModePerm)
	if err != nil {
		resp["err"] = "unable to proceed"
		handler.logger.Error("err creating public directory", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}

	// Full path to the new public location
	newPath := filepath.Join(publicDir, era5_response.FileName)
	err = os.Rename(era5_response.FileName, newPath)
	if err != nil {
		resp["err"] = "unable to proceed"
		handler.logger.Error("err moving file to public directory", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}

	era5_response.FileURL = fmt.Sprintf("http://localhost:8080/files/%s", era5_response.FileName)
	resp["message"] = "Data fetched successfully"
	apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusOK)
	handler.logger.Info("Data fetched Successfully")
}
