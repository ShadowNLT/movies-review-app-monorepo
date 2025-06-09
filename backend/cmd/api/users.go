package main

import (
	"cinepulse.nlt.net/internal/data/users"
	"cinepulse.nlt.net/internal/data/users/inputs"
	"cinepulse.nlt.net/internal/validator"
	"errors"
	"net/http"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input inputs.CreateUserInput

	err := app.readJSON(w, r, &input, 2048)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext
	err = input.Password.Set(input.PasswordString)

	v := validator.New()

	if err != nil {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Validate the input struct
	if inputs.ValidateCreateUserInput(v, &input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.Insert(&input)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrDuplicateEmail):
			v.AddError("email", "a user with this Email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, users.ErrDuplicateProfileHandle):
			v.AddError("profile_handle", "a user with this ProfileHandle already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
