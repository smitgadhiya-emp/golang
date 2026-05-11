package middleware

import (
	"errors"
	"gin-project/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))

	if token == "" {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", errors.New("missing authorization token"))
		c.Abort()
		return
	}

	if !helper.ValidateJWT(token) {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", errors.New("invalid token"))
		c.Abort()
		return
	}

	claims, err := helper.GetJWTClaims(token)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", errors.New("invalid token claims"))
		c.Abort()
		return
	}

	userID, ok := claims["userId"].(string)
	if !ok || userID == "" {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", errors.New("missing user id in token"))
		c.Abort()
		return
	}

	c.Set("userID", userID)
}
