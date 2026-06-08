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
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

// resetTokenTTL is how long a password-reset link stays valid.
const resetTokenTTL = 15 * time.Minute

// resetTokenKey builds the Redis key under which a reset token is stored.
func resetTokenKey(token string) string {
	return "reset_password:" + token
}

// resetLinkBase returns the deep-link base the reset email points at.
func resetLinkBase() string {
	if base := os.Getenv("RESET_PASSWORD_DEEPLINK"); base != "" {
		return base
	}
	return "myayur://reset-password"
}

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

// ForgotPassword creates a one-time reset token, stores it in Redis against the
// user's id, and emails a deep link the user can tap to reset their password.
func ForgotPassword(data dto.ForgotPasswordPayload) (*dto.ForgotPasswordResponse, error) {
	user, err := repositories.GetUserByEmail(data.Email)
	if err != nil {
		return nil, fmt.Errorf("user does not exist with this email")
	}

	token, err := helper.GenerateRandomToken()
	if err != nil {
		return nil, err
	}

	// Store token -> userID with a short TTL.
	if err := config.Redis.Set(config.Ctx, resetTokenKey(token), user.ID, resetTokenTTL).Err(); err != nil {
		return nil, fmt.Errorf("failed to create reset token")
	}

	resetLink := resetLinkBase() + "?token=" + token

	// Send the reset email asynchronously via the email queue.
	if err := queue.PublishResetPasswordEmail(user.Email, user.UserName, resetLink, int(resetTokenTTL.Minutes())); err != nil {
		log.Printf("failed to enqueue reset-password email for %s: %v", user.Email, err)
	}

	return &dto.ForgotPasswordResponse{
		Message: "If an account exists for this email, a reset link has been sent",
	}, nil
}

// VerifyResetToken checks that a reset token is still valid (exists in Redis)
// without consuming it, so the app can decide whether to show the reset form.
func VerifyResetToken(data dto.VerifyResetTokenPayload) (*dto.VerifyResetTokenResponse, error) {
	if data.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	userID, err := config.Redis.Get(config.Ctx, resetTokenKey(data.Token)).Result()
	if err != nil {
		return nil, fmt.Errorf("invalid or expired reset token")
	}

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &dto.VerifyResetTokenResponse{
		Valid: true,
		Email: user.Email,
	}, nil
}

// ResetPassword validates the reset token, sets the new password, and deletes
// the token so it can't be reused.
func ResetPassword(data dto.ResetPasswordPayload) (*dto.ResetPasswordResponse, error) {
	if data.Token == "" {
		return nil, fmt.Errorf("token is required")
	}
	if len(data.NewPassword) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}

	key := resetTokenKey(data.Token)

	userID, err := config.Redis.Get(config.Ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("invalid or expired reset token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if err := repositories.UpdateUserPassword(userID, string(hashedPassword)); err != nil {
		return nil, err
	}

	// One-time use: invalidate the token now that the password is changed.
	if err := config.Redis.Del(config.Ctx, key).Err(); err != nil {
		log.Printf("failed to delete used reset token: %v", err)
	}

	return &dto.ResetPasswordResponse{
		Message: "Password has been reset successfully",
	}, nil
}
