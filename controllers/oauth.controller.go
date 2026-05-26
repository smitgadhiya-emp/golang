package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"

	"gin-project/config"
	"gin-project/services"

	"github.com/gin-gonic/gin"
)

const oauthStateCookie = "oauth_state"

func GoogleOAuthStart(c *gin.Context) {
	oauthConfig, err := config.GoogleOAuthConfig()
	if err != nil {
		redirectOAuthError(c, "Google sign-in is not configured on the server")
		return
	}

	state, err := randomState()
	if err != nil {
		redirectOAuthError(c, "Could not start Google sign-in")
		return
	}

	c.SetCookie(oauthStateCookie, state, 600, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, oauthConfig.AuthCodeURL(state))
}

func GoogleOAuthCallback(c *gin.Context) {
	if errMsg := c.Query("error"); errMsg != "" {
		redirectOAuthError(c, "Google sign-in was cancelled")
		return
	}

	state := c.Query("state")
	cookieState, err := c.Cookie(oauthStateCookie)
	if err != nil || state == "" || state != cookieState {
		redirectOAuthError(c, "Invalid OAuth state")
		return
	}

	c.SetCookie(oauthStateCookie, "", -1, "/", "", false, true)

	oauthConfig, err := config.GoogleOAuthConfig()
	if err != nil {
		redirectOAuthError(c, "Google sign-in is not configured on the server")
		return
	}

	token, err := oauthConfig.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		redirectOAuthError(c, "Failed to exchange Google authorization code")
		return
	}

	googleUser, err := config.FetchGoogleUserInfo(c.Request.Context(), token)
	if err != nil {
		redirectOAuthError(c, "Failed to fetch Google profile")
		return
	}

	login, err := services.LoginOrRegisterWithGoogle(googleUser)
	if err != nil {
		redirectOAuthError(c, "Failed to sign in with Google")
		return
	}

	params := url.Values{}
	params.Set("token", login.Token)
	params.Set("userId", login.UserID)
	params.Set("email", login.Email)
	params.Set("role", login.Role)

	c.Redirect(http.StatusTemporaryRedirect, config.FrontendURL()+"/oauth/callback?"+params.Encode())
}

func redirectOAuthError(c *gin.Context, message string) {
	params := url.Values{}
	params.Set("error", message)
	c.Redirect(http.StatusTemporaryRedirect, config.FrontendURL()+"/oauth/callback?"+params.Encode())
}

func randomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
