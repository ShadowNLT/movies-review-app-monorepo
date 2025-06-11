package types

type UserWelcomeTemplateData struct {
	ProfileHandle  string `json:"profileHandle"`
	CurrentYear    int    `json:"currentYear"`
	ActivationLink string `json:"activationLink"`
}
