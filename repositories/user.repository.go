package repositories

import (
	"gin-project/config"
	"gin-project/dto"
	"gin-project/entity"
)

func CreateUser(data dto.SignupPayload) (string, error) {
	user := entity.User{
		UserName: data.UserName,
		Password: data.Password,
		Email:    data.Email,
		City:     data.City,
		Pincode:  data.Pincode,
		Role:     data.Role,
		Provider: "local",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return "", err
	}

	return user.ID, nil
}

func GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(userID string) (*entity.User, error) {
	var user entity.User

	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserPassword(userID string, password string) error {
	return config.DB.Model(&entity.User{}).Where("id = ?", userID).Update("password", password).Error
}

func GetUserByGoogleID(googleID string) (*entity.User, error) {
	var user entity.User
	if err := config.DB.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateGoogleUser(userName, email, googleID, passwordHash string) (*entity.User, error) {
	user := entity.User{
		UserName: userName,
		Password: passwordHash,
		Email:    email,
		City:     "Unknown",
		Pincode:  0,
		Role:     "user",
		Provider: "google",
		GoogleID: googleID,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func LinkGoogleAccount(userID, googleID string) error {
	return config.DB.Model(&entity.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"google_id": googleID,
		"provider":  "google",
	}).Error
}
