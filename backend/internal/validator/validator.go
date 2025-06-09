package validator

import (
	"regexp"
	"slices"
)

// EmailRX regular expression pattern is
// taken from https://html.spec.whatwg.org/#valid-e-mail-address.
var (
	EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
)

type ErrorMap = map[string]string
type Validator struct {
	Errors ErrorMap
}

// New returns a new Validator Instance with an empty errors map
func New() *Validator {
	return &Validator{Errors: make(ErrorMap)}
}

// Valid returns true if the errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(errorKey, errorMessage string) {
	if _, exists := v.Errors[errorKey]; !exists {
		v.Errors[errorKey] = errorMessage
	}
}

// AddErrorIfNot adds a new error entry to the errors map if the predicate is equal to true
func (v *Validator) AddErrorIfNot(predicate bool, errorKey, errorMessage string) {
	if !predicate {
		v.AddError(errorKey, errorMessage)
	}
}

func (v *Validator) RequiredString(value, errorKey string) {
	v.AddErrorIfNot(value != "", errorKey, "must be provided")
}

// PermittedValue returns true if a value is in a list of permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
