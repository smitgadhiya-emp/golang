package services

import (
	"fmt"
	"gin-project/dto"
	"gin-project/helper"
	"gin-project/repositories"
	"golang.org/x/crypto/bcrypt"
)

func SignUpService(data dto.SignupPayload) (*dto.SignupResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	data.Password = string(hashedPassword)

	userID, err := repositories.CreateUser(data)
	if err != nil {
		return nil, err
	}

	token, err := helper.GenerateJWT(userID, data.Email, data.Role)
	if err != nil {
		return nil, err
	}

	return &dto.SignupResponse{
		ID:    userID,
		Token: token,
	}, nil
}

func Login(data dto.LoginPayload) (*dto.LoginResponse, error) {

	email := data.Email
	password := data.Password

	// find user by email
	getUser, err := repositories.GetUserByEmail(email)

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// validate password
	isValid := helper.ValidatePassword(getUser.Password, password)

	if !isValid {
		return nil, fmt.Errorf("invalid credentials")
	}

	// generate JWT token
	token, err := helper.GenerateJWT(getUser.ID, getUser.Email, getUser.Role)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:  token,
		UserID: getUser.ID,
		Email:  getUser.Email,
		Role:   getUser.Role,
	}, nil

}

func ChangePassword(userID string, data dto.ChangePasswordPayload) (*dto.ChangePasswordResponse, error) {
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !helper.ValidatePassword(user.Password, data.OldPassword) {
		return nil, fmt.Errorf("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if err := repositories.UpdateUserPassword(userID, string(hashedPassword)); err != nil {
		return nil, err
	}

	return &dto.ChangePasswordResponse{
		Message: "Password changed successfully",
	}, nil
}
