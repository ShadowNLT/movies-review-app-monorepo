package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		// There is a Syntax Error with the "JSON" Body sent over the request
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// This happens when the JSON value is the wrong type for the
		// target destination
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)

			}

		// This happens if the request body is empty
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// This happens if we pass something that is not a non-nil pointer as destination
		// which should never be the case to begin with. This is why it makes sense to panic
		// here as this error only happen due to a mistake in the code
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}
