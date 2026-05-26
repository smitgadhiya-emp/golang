package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func GoogleOAuthConfig() (*oauth2.Config, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("OAUTH_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/api/v1/oauth/google/callback"
	}

	if clientID == "" || clientSecret == "" || clientID == "YOUR_CLIENT_ID" {
		return nil, fmt.Errorf("google oauth is not configured")
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}, nil
}

func FetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned status %d", res.StatusCode)
	}

	var info GoogleUserInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, err
	}

	if info.Email == "" {
		return nil, fmt.Errorf("google account has no email")
	}

	return &info, nil
}

func FrontendURL() string {
	if url := os.Getenv("FRONTEND_URL"); url != "" {
		return url
	}
	return "http://localhost:3000"
}
