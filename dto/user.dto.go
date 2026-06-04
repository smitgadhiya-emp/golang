package dto

import "time"

type UpdateUserPayload struct {
	UserName string `json:"userName"`
	Phone    string `json:"phone"`
	City     string `json:"city"`
	Pincode  int    `json:"pincode"`
	Role     string `json:"role"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	UserName  string    `json:"userName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	City      string    `json:"city"`
	Pincode   int       `json:"pincode"`
	Role      string    `json:"role"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"createdAt"`
}
