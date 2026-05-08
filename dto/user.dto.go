package dto

type SignupPayload struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
	City     string `json:"city"`
	Pincode  int    `json:"pincode"`
	Role     string `json:"role"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
