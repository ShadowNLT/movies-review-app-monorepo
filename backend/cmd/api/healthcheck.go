package main

import (
	"cinepulse.nlt.net/internal/constants"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     appVersion,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, constants.ErrorMessages[http.StatusInternalServerError], http.StatusInternalServerError)
		return
	}
}
