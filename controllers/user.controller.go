package controllers

import (
	"errors"
	"gin-project/dto"
	"gin-project/helper"
	"gin-project/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", errors.New("missing user id"))
		return
	}

	res, err := services.GetMe(userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "Profile retrieved successfully", res)
}

func GetAllUsers(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		helper.ErrorResponse(c, http.StatusForbidden, "Forbidden", errors.New("admin access required"))
		return
	}

	res, err := services.GetAllUsers()
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "Users retrieved successfully", res)
}

func GetUserByID(c *gin.Context) {
	callerID := c.GetString("userID")
	role := c.GetString("role")
	targetID := c.Param("id")

	if role != "admin" && callerID != targetID {
		helper.ErrorResponse(c, http.StatusForbidden, "Forbidden", errors.New("cannot view other user's profiles"))
		return
	}

	res, err := services.GetUserByID(targetID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "User retrieved successfully", res)
}

func UpdateUser(c *gin.Context) {
	callerID := c.GetString("userID")
	role := c.GetString("role")
	targetID := c.Param("id")

	if role != "admin" && callerID != targetID {
		helper.ErrorResponse(c, http.StatusForbidden, "Forbidden", errors.New("cannot update other user's profiles"))
		return
	}

	var req dto.UpdateUserPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Only admin can change user roles
	if req.Role != "" && role != "admin" {
		helper.ErrorResponse(c, http.StatusForbidden, "Forbidden", errors.New("only admins can modify roles"))
		return
	}

	res, err := services.UpdateUser(targetID, req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to update user", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "User updated successfully", res)
}

func DeleteUser(c *gin.Context) {
	callerID := c.GetString("userID")
	role := c.GetString("role")
	targetID := c.Param("id")

	if role != "admin" && callerID != targetID {
		helper.ErrorResponse(c, http.StatusForbidden, "Forbidden", errors.New("cannot delete other user's accounts"))
		return
	}

	err := services.DeleteUser(targetID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to delete user", err)
		return
	}

	helper.SucessResponse(c, http.StatusOK, "User deleted successfully", gin.H{"id": targetID})
}
