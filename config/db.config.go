package config

import (
	"database/sql"
	"gin-project/entity"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB() {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	DB = db
	entity.MigrateUserTable(DB)
	log.Println("SQLite connected")
}
