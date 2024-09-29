package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/h3th-IV/climateer/pkg/database"
	"github.com/h3th-IV/climateer/pkg/model"
	"github.com/h3th-IV/climateer/pkg/utils"
	"go.uber.org/zap"
)

type measurementHandler struct {
	logger      *zap.Logger
	mysqlclient database.Database
}

func NewMeasurementHandler(logger *zap.Logger, mysqlclient database.Database) *measurementHandler {
	return &measurementHandler{
		logger:      logger,
		mysqlclient: mysqlclient,
	}
}

func (handler *measurementHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var TTL = 60
	var resp = map[string]interface{}{}
	var measurement model.UserMeasurement

	// Authenticate the user
	userInfo, err := utils.AuthenticateUser(r.Context(), handler.logger, handler.mysqlclient)
	if err != nil {
		resp["err"] = "please sign in to submit measurements"
		handler.logger.Warn("unauthorized user")
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusUnauthorized)
		return
	}

	// Decode the incoming request body into measurement struct
	if err := json.NewDecoder(r.Body).Decode(&measurement); err != nil {
		resp["err"] = "please provide values to your measurement"
		handler.logger.Error("error decoding request body", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusBadRequest)
		return
	}

	// Fetch the country ID based on the country name
	countryID, err := handler.mysqlclient.GetCountryIDByName(r.Context(), measurement.CountryName)
	if err != nil || countryID == 0 {
		resp["err"] = "country not found"
		handler.logger.Error("error fetching country ID", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusBadRequest)
		return
	}

	indicatorID, err := handler.mysqlclient.GetIndicatorIDByCode(r.Context(), measurement.IndicatorCode)
	if err != nil || indicatorID == 0 {
		resp["err"] = "indicator not found"
		handler.logger.Error("error fetching indicator ID", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusBadRequest)
		return
	}
	// Insert user measurement into the database
	success, err := handler.mysqlclient.AddUserMeasurement(
		r.Context(),
		userInfo.ID,       // User ID
		countryID,         // Country ID
		indicatorID,       // Indicator ID
		measurement.Year,  // Year of measurement
		measurement.Value, // Value of the measurement
	)

	if err != nil || !success {
		resp["err"] = "unable to submit measurement"
		handler.logger.Error("error inserting user measurement", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}

	// Success response
	resp["message"] = "measurement successfully submitted"
	apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusOK)
	handler.logger.Info("user measurement submitted successfully")
}
