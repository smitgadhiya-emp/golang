package repositories

import (
	"gin-project/config"
	"gin-project/dto"
)

func CreateUser(data dto.SignupPayload) (string, error) {
	query := `
	INSERT INTO users (
		id,
		userName,
		password,
		email,
		city,
		pincode,
		role
	)
	VALUES (
		lower(hex(randomblob(16))),
		?,
		?,
		?,
		?,
		?,
		?
	)
	RETURNING id;`

	var userID string
	err := config.DB.QueryRow(
		query,
		data.UserName,
		data.Password,
		data.Email,
		data.City,
		data.Pincode,
		data.Role,
	).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
