package controllers

import (
	"net/http"

	"gin-project/dto"
	"gin-project/helper"

	"github.com/gin-gonic/gin"
)

func SendWelcomeEmail(c *gin.Context) {
	var req dto.SendWelcomeEmailPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	if err := helper.SendWelcomeEmail(req.Email, req.Name); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to send welcome email", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "Welcome email sent successfully", nil)
}

func SendOtpEmail(c *gin.Context) {
	var req dto.SendOtpEmailPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	if err := helper.SendOTPEmail(req.Email, req.Name, req.Otp, req.ExpiryMinutes); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to send OTP email", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "OTP email sent successfully", nil)
}
