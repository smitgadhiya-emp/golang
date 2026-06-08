package services

import (
	"fmt"
	"gin-project/config"
	"gin-project/dto"
	"gin-project/entity"
	"gin-project/helper"
	"gin-project/queue"
	"gin-project/repositories"
	"log"

	"golang.org/x/crypto/bcrypt"

	verify "github.com/twilio/twilio-go/rest/verify/v2"
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

	// Enqueue the welcome email; the worker sends it asynchronously so signup
	// stays fast and isn't blocked on SMTP. A publish failure must not fail
	// registration, so we only log it.
	if err := queue.PublishWelcomeEmail(data.Email, data.UserName); err != nil {
		log.Printf("failed to enqueue welcome email for %s: %v", data.Email, err)
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

func SendOTP(data dto.SendOtpPayload) (bool, error) {

	phoneNumber := "+" + data.PhoneCountryCode + data.PhoneNumber

	client := config.GetTwilioClient()

	params := &verify.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(
		config.GetTwilioVerifyServiceSID(),
		params,
	)
	fmt.Printf("resp %+v", resp)

	if err != nil {
		return false, err
	}

	fmt.Println(resp.Status)
	return true, nil
}

func VerifyOTP(data dto.VerifyOtpPayload, userDetails *entity.User) (*dto.LoginResponse, error) {

	phone := "+" + data.PhoneCountryCode + data.PhoneNumber
	otp := data.Otp

	client := config.GetTwilioClient()

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(otp)

	resp, err := client.VerifyV2.CreateVerificationCheck(
		config.GetTwilioVerifyServiceSID(),
		params,
	)

	if err != nil {
		return nil, err
	}

	if resp.Status == nil || *resp.Status != "approved" {
		return nil, fmt.Errorf("invalid or expired otp")
	}

	getUser := userDetails

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
