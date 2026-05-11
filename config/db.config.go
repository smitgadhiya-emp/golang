package config

import (
	"gin-project/entity"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("./db/app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	entity.MigrateUserTable(DB)
	log.Println("SQLite connected")
}
