package shared

import (
	"cinepulse.nlt.net/internal/validator"
	"time"
	"unicode/utf8"
)

func ValidateEmail(v *validator.Validator, email string) {
	v.RequiredString(email, "email")
	v.AddErrorIfNot(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidateProfileHandle(v *validator.Validator, handle string) {
	v.RequiredString(handle, "handle")
	v.AddErrorIfNot(utf8.RuneCountInString(handle) >= 3, "handle", "must be at least 3 characters")
	v.AddErrorIfNot(utf8.RuneCountInString(handle) <= 30, "handle", "must not be more than 30 characters")
}

func IsAtLeast13YearsOld(birthDate time.Time) bool {
	limit := time.Now().AddDate(-13, 0, 0)
	return birthDate.Before(limit) || birthDate.Equal(limit)
}

func IsAtMost100YearsOld(birthDate time.Time) bool {
	limit := time.Now().AddDate(-100, 0, 0)
	return birthDate.After(limit) || birthDate.Equal(limit)
}

func ValidatePasswordPlaintext(v *validator.Validator, pwd string) {
	v.RequiredString(pwd, "Password")
	v.AddErrorIfNot(len(pwd) >= 8, "Password", "must be at least 8 bytes long")
	v.AddErrorIfNot(len(pwd) <= 72, "Password", "must not be more than 72 bytes long")
}
