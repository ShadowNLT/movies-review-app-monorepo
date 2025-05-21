package main

import (
	"cinepulse.nlt.net/internal/constants"
	"fmt"
	"net/http"
)

// Generic helper for logging an error message
// along with the current request's method and URL
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.RequestURI
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// Generic helper for sending JSON-formatted error messages to the client with
// a given status code
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// This helper is for when we encounter an unexpected problem. A problem that shouldn't happen
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusInternalServerError, constants.ErrorMessages[http.StatusInternalServerError])
}

// This helper is for when we want to send a 404 to the client
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, constants.ErrorMessages[http.StatusNotFound])
}

// This helper is for when we want to send a 405 Method Not Allowed to the client
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf(constants.FormattedErrorMessages[http.StatusMethodNotAllowed], r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// This helper is for when we want to send a 400 Bad request to the client
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
