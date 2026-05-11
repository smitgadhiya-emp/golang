package dto

type SignupPayload struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
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

type ChangePasswordPayload struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}
