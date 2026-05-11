package entity

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	UserName  string    `gorm:"column:userName;not null" json:"userName"`
	Password  string    `gorm:"column:password;not null" json:"-"`
	Email     string    `gorm:"column:email;unique;not null" json:"email"`
	City      string    `gorm:"column:city;not null" json:"city"`
	Pincode   int       `gorm:"column:pincode;not null" json:"pincode"`
	Role      string    `gorm:"column:role;not null" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (User) TableName() string {
	return "users"
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.ID != "" {
		return nil
	}

	id, err := generateID()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func MigrateUserTable(db *gorm.DB) {
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal(err)
	}

	seedDefaultUser(db)
}

func seedDefaultUser(db *gorm.DB) {
	defaultEmail := "admin@example.com"

	var user User
	err := db.Where(&User{Email: defaultEmail}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		user = User{
			UserName: "admin",
			Password: string(hashedPassword),
			Email:    defaultEmail,
			City:     "Mumbai",
			Pincode:  400001,
			Role:     "admin",
		}

		if err := db.Create(&user).Error; err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	} else if user.Password == "admin123" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		if err := db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
			log.Fatal(err)
		}
	}

	log.Println("\n ========================= \n Default user ensured: \n email: admin@example.com \n Password: admin123 \n ========================= \n")
}

func generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
