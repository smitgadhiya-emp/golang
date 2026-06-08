package helper

import (
	"os"
	"time"

	"gin-project/config"
	"gin-project/templates"
)

const defaultAppName = "MyAyur"

// appName returns the configurable product name used across email templates.
func appName() string {
	if name := os.Getenv("APP_NAME"); name != "" {
		return name
	}
	return defaultAppName
}

// appURL returns the link used by call-to-action buttons in emails.
func appURL() string {
	if url := os.Getenv("FRONTEND_URL"); url != "" {
		return url
	}
	return "https://golang-1-w8al.onrender.com"
}

type welcomeEmailData struct {
	Name    string
	AppName string
	AppUrl  string
	Year    int
}

type otpEmailData struct {
	Name          string
	OTP           string
	ExpiryMinutes int
	AppName       string
	Year          int
}

type resetPasswordEmailData struct {
	Name          string
	ResetLink     string
	ExpiryMinutes int
	AppName       string
	Year          int
}

// SendWelcomeEmail renders the welcome template and emails it to the user.
func SendWelcomeEmail(to, name string) error {
	body, err := templates.Render("welcome.html", welcomeEmailData{
		Name:    name,
		AppName: appName(),
		AppUrl:  appURL(),
		Year:    time.Now().Year(),
	})
	if err != nil {
		return err
	}

	return config.SendMail(to, "Welcome to "+appName()+" 🌿", body)
}

// SendOTPEmail renders the OTP template and emails the verification code.
func SendOTPEmail(to, name, otp string, expiryMinutes int) error {
	if expiryMinutes <= 0 {
		expiryMinutes = 10
	}

	body, err := templates.Render("otp.html", otpEmailData{
		Name:          name,
		OTP:           otp,
		ExpiryMinutes: expiryMinutes,
		AppName:       appName(),
		Year:          time.Now().Year(),
	})
	if err != nil {
		return err
	}

	return config.SendMail(to, "Your "+appName()+" verification code", body)
}

// SendResetPasswordEmail renders the reset-password template and emails the
// reset link to the user.
func SendResetPasswordEmail(to, name, resetLink string, expiryMinutes int) error {
	if expiryMinutes <= 0 {
		expiryMinutes = 15
	}

	body, err := templates.Render("reset-password.html", resetPasswordEmailData{
		Name:          name,
		ResetLink:     resetLink,
		ExpiryMinutes: expiryMinutes,
		AppName:       appName(),
		Year:          time.Now().Year(),
	})
	if err != nil {
		return err
	}

	return config.SendMail(to, "Reset your "+appName()+" password", body)
}
