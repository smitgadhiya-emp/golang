package controllers

import (
	"fmt"
	"gin-project/dto"
	"gin-project/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Singup(c *gin.Context) {

	var req dto.SignupPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := services.SignUpService(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"id":      res,
	})
}

func Login(c *gin.Context) {

	var req dto.LoginPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("req body data %+v \n", req)

	c.JSON(200, gin.H{
		"message": "login",
	})
}
