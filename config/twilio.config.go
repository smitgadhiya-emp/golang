package config

import (
	"os"
	"sync"

	"github.com/twilio/twilio-go"
)

var (
	TwilioClient *twilio.RestClient
	once         sync.Once
)

// InitTwilio initializes the Twilio RestClient.
func InitTwilio() {
	once.Do(func() {
		accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
		authToken := os.Getenv("TWILIO_AUTH_TOKEN")
		if accountSid == "" || authToken == "" {
			// Fall back to automatic environment resolution
			TwilioClient = twilio.NewRestClient()
		} else {
			TwilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
				Username: accountSid,
				Password: authToken,
			})
		}
	})
}

// GetTwilioClient returns the global Twilio RestClient instance.
func GetTwilioClient() *twilio.RestClient {
	if TwilioClient == nil {
		InitTwilio()
	}
	return TwilioClient
}

// GetTwilioVerifyServiceSID retrieves the Twilio Verify Service SID.
func GetTwilioVerifyServiceSID() string {
	return os.Getenv("TWILIO_VERIFY_SERVICE_SID")
}
