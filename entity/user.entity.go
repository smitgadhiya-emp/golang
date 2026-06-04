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
	Phone     string    `gorm:"column:phone;unique;not null" json:"phone"`
	City      string    `gorm:"column:city;not null" json:"city"`
	Pincode   int       `gorm:"column:pincode;not null" json:"pincode"`
	Role      string    `gorm:"column:role;not null" json:"role"`
	Provider  string    `gorm:"column:provider;not null;default:local" json:"provider"`
	GoogleID  *string   `gorm:"column:google_id;uniqueIndex" json:"-"`
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
	backfillPhones(db)

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal(err)
	}

	seedDefaultUser(db)
}

// backfillPhones gives every existing user a unique, non-empty phone before
// AutoMigrate builds the UNIQUE index on the column. Without this, older rows
// that were created with the empty-string default collide on the new
// constraint ("UNIQUE constraint failed: users.phone").
func backfillPhones(db *gorm.DB) {
	migrator := db.Migrator()

	// A previously failed table rebuild can leave GORM's temp table behind,
	// which would break the next AutoMigrate ("table users__temp already
	// exists"). Drop it if present.
	if err := db.Exec("DROP TABLE IF EXISTS users__temp").Error; err != nil {
		log.Fatal(err)
	}

	// Fresh database: the table doesn't exist yet, so AutoMigrate will create
	// it cleanly with no rows to back-fill.
	if !migrator.HasTable(&User{}) {
		return
	}

	// Add the column without any constraint first so we can populate it before
	// the UNIQUE index is created.
	if !migrator.HasColumn(&User{}, "phone") {
		if err := db.Exec("ALTER TABLE users ADD COLUMN phone text DEFAULT ''").Error; err != nil {
			log.Fatal(err)
		}
	}

	// Derive a unique placeholder from each user's (unique) id. Users can set a
	// real number later from the profile screen.
	if err := db.Exec("UPDATE users SET phone = '+0' || id WHERE phone IS NULL OR phone = ''").Error; err != nil {
		log.Fatal(err)
	}

	// Non-Google users must store NULL (not '') for google_id, otherwise the
	// second local user collides on the unique index.
	if migrator.HasColumn(&User{}, "google_id") {
		if err := db.Exec("UPDATE users SET google_id = NULL WHERE google_id = ''").Error; err != nil {
			log.Fatal(err)
		}
	}
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
			Phone:    "+910000000001",
			City:     "Mumbai",
			Pincode:  400001,
			Role:     "admin",
			Provider: "local",
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
