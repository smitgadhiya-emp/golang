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
