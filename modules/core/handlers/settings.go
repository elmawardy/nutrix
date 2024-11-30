package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
)

// UpdateSettings is a post request handler that updates the settings in the database
// send a models.Settings directory to body to use it.
func UpdateSettings(conf config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var settings models.Settings
		err := json.NewDecoder(r.Body).Decode(&settings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		settings_svc := services.SettingsService{
			Config: conf,
		}

		err = settings_svc.UpdateSettings(settings)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// GetSettings is an http get handlers that just returns the settings object from the db
func GetSettings(conf config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		settings_svc := services.SettingsService{
			Config: conf,
		}

		settings, err := settings_svc.GetSettings()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			Body models.Settings `json:"settings"`
		}{
			Body: settings,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to marshal order settings response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}

}
