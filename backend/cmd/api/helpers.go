package main

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, statusCode int, data envelope, headers http.Header) error {
	// Format the JSON to make it easier to read on terminal apps
	jsonBytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Add a new line at the end to make it easier to view in terminal applications
	jsonBytes = append(jsonBytes, '\n')

	// Go through the header map and set the headers that we need to have set
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonBytes)
	if err != nil {
		return err
	}
	return nil
}
