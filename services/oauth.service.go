package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"gin-project/config"
	"gin-project/dto"
	"gin-project/entity"
	"gin-project/helper"
	"gin-project/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginOrRegisterWithGoogle(info *config.GoogleUserInfo) (*dto.LoginResponse, error) {
	user, err := repositories.GetUserByGoogleID(info.ID)
	if err == nil {
		return issueLoginResponse(user)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	existing, err := repositories.GetUserByEmail(info.Email)
	if err == nil {
		if existing.GoogleID == "" {
			if linkErr := repositories.LinkGoogleAccount(existing.ID, info.ID); linkErr != nil {
				return nil, linkErr
			}
			existing.GoogleID = info.ID
			existing.Provider = "google"
		}
		return issueLoginResponse(existing)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	userName := info.Name
	if userName == "" {
		userName = info.Email
	}

	passwordHash, err := randomPasswordHash()
	if err != nil {
		return nil, err
	}

	user, err = repositories.CreateGoogleUser(userName, info.Email, info.ID, passwordHash)
	if err != nil {
		return nil, err
	}

	return issueLoginResponse(user)
}

func issueLoginResponse(user *entity.User) (*dto.LoginResponse, error) {
	token, err := helper.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{
		Token:  token,
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

func randomPasswordHash() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(hex.EncodeToString(bytes)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
