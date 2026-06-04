package dto

type SignupPayload struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	City     string `json:"city"`
	Pincode  int    `json:"pincode"`
	Role     string `json:"role"`
}

type SignupResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type SendOtpPayload struct {
	PhoneCountryCode string `json:"phoneCountryCode"`
	PhoneNumber      string `json:"phoneNumber"`
}
type VerifyOtpPayload struct {
	PhoneCountryCode string `json:"phoneCountryCode"`
	PhoneNumber      string `json:"phone"`
	Otp              string `json:"otp"`
}

type ChangePasswordPayload struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}
