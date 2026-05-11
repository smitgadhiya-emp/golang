package controllers

import (
	"gin-project/dto"
	"gin-project/helper"
	"gin-project/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Singup(c *gin.Context) {

	var req dto.SignupPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
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
