package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     appVersion,
	}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Add a new line at the end to make it easier to view in terminal applications
	jsonValue = append(jsonValue, '\n')

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonValue)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
