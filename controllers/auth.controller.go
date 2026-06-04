package controllers

import (
	"fmt"
	"gin-project/dto"
	"gin-project/helper"
	"gin-project/repositories"
	"gin-project/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Singup(c *gin.Context) {

	var req dto.SignupPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	res, err := services.SignUpService(req)

	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to create user", err)
		return
	}

	helper.SucessResponse(c, http.StatusCreated, "User created successfully", res)

}

func Login(c *gin.Context) {

	var req dto.LoginPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	res, err := services.Login(req)

	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid credentials", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "Login successful", res)
}

func SendOTP(c *gin.Context) {

	var req dto.SendOtpPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// find the user registered with this phone number

	phone := "+" + req.PhoneCountryCode + req.PhoneNumber
	_, err := repositories.GetUserByPhone(phone)

	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "User is not exist with this phone number", err)
		return
	}

	res, err := services.SendOTP(req)

	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "Otp send successfully", res)
}

func VerifyOTP(c *gin.Context) {
	var req dto.VerifyOtpPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// find the user registered with this phone number

	phone := "+" + req.PhoneCountryCode + req.PhoneNumber

	getUser, err := repositories.GetUserByPhone(phone)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "User is not exist with this phone number", err)
		return
	}

	res, err := services.VerifyOTP(req, getUser)

	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "Otp verified successfully", res)
}

func ChangePassword(c *gin.Context) {

	var req dto.ChangePasswordPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("missing user id"))
		return
	}

	res, err := services.ChangePassword(userID, req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to change password", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "Password changed successfully", res)

}
