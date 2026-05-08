package services

import (
	"gin-project/dto"
	"gin-project/repositories"
)

func SignUpService(data dto.SignupPayload) (string, error) {

	return repositories.CreateUser(data)

}
