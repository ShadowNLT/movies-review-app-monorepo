package inputs

import (
	"cinepulse.nlt.net/internal/data/users/shared"
	"cinepulse.nlt.net/internal/validator"
	"time"
)

type CreateUserInput struct {
	Email          string          `json:"email"`
	PasswordString string          `json:"password"`
	Password       shared.Password `json:"-"`
	ProfileHandle  string          `json:"profile_handle"`
	Location       string          `json:"location"`
	DateOfBirth    time.Time       `json:"date_of_birth"`
}

func ValidateCreateUserInput(v *validator.Validator, u *CreateUserInput) {
	// ProfileHandle
	shared.ValidateProfileHandle(v, u.ProfileHandle)

	// E-mail
	shared.ValidateEmail(v, u.Email)

	// Password
	if u.Password.Plaintext != nil {
		shared.ValidatePasswordPlaintext(v, *u.Password.Plaintext)
	}

	// IF the Password hash is ever nil, we need to panic as this should never happen
	if u.Password.Hash == nil {
		panic("missing Password hash for user")
	}

	// Location
	v.RequiredString(u.Location, "location")

	// Date of birth
	v.AddErrorIfNot(shared.IsAtLeast13YearsOld(u.DateOfBirth), "date_of_birth", "must be at least 13 years old")
	v.AddErrorIfNot(shared.IsAtMost100YearsOld(u.DateOfBirth), "date_of_birth", "must be at most 100 years old")
}
