package dto

type SendWelcomeEmailPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SendOtpEmailPayload struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Otp           string `json:"otp"`
	ExpiryMinutes int    `json:"expiryMinutes"`
}
