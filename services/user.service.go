package services

import (
	"gin-project/dto"
	"gin-project/entity"
	"gin-project/repositories"
)

func toUserResponse(user *entity.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	return &dto.UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		Phone:     user.Phone,
		City:      user.City,
		Pincode:   user.Pincode,
		Role:      user.Role,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
	}
}

func GetMe(userID string) (*dto.UserResponse, error) {
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

func GetAllUsers() ([]dto.UserResponse, error) {
	users, err := repositories.GetAllUsers()
	if err != nil {
		return nil, err
	}

	response := make([]dto.UserResponse, len(users))
	for i, u := range users {
		response[i] = *toUserResponse(&u)
	}
	return response, nil
}

func GetUserByID(userID string) (*dto.UserResponse, error) {
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

func UpdateUser(userID string, payload dto.UpdateUserPayload) (*dto.UserResponse, error) {
	user, err := repositories.UpdateUser(userID, payload)
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

func DeleteUser(userID string) error {
	return repositories.DeleteUser(userID)
}
