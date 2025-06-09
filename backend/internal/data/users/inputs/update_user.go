package inputs

import "cinepulse.nlt.net/internal/validator"

type UpdateUserInput struct {
	Email         *string `json:"email"`
	ProfileHandle *string `json:"profile_handle"`
	Location      *string `json:"location"`
	IsProtected   *bool   `json:"is_protected"` // Does not cause version change
}

func ValidateUpdateUserInput(v *validator.Validator, input *UpdateUserInput) {
	if input.Email == nil && input.ProfileHandle == nil && input.Location == nil && input.IsProtected == nil {
		v.AddError("all", "at least one of email, profile_handle, location, is_protected must be set")
	}
}
