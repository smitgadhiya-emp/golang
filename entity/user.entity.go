package entity

import (
	"database/sql"
	"log"
)

func MigrateUserTable(db *sql.DB) {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		userName TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		city TEXT NOT NULL,
		pincode TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createUsersTable); err != nil {
		log.Fatal(err)
	}

	ensureUserIDIsGeneratedText(db)
	seedDefaultUser(db)
}

func ensureUserIDIsGeneratedText(db *sql.DB) {
	var idType string
	err := db.QueryRow(`
		SELECT type
		FROM pragma_table_info('users')
		WHERE name = 'id';
	`).Scan(&idType)
	if err != nil {
		log.Fatal(err)
	}

	if idType == "TEXT" {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`
		CREATE TABLE users_new (
			id TEXT PRIMARY KEY,
			userName TEXT NOT NULL,
			password TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			city TEXT NOT NULL,
			pincode TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`); err != nil {
		log.Fatal(err)
	}

	if _, err = tx.Exec(`
		INSERT INTO users_new (
			id,
			userName,
			password,
			email,
			city,
			pincode,
			role,
			created_at
		)
		SELECT
			lower(hex(randomblob(16))),
			userName,
			password,
			email,
			city,
			pincode,
			role,
			created_at
		FROM users;
	`); err != nil {
		log.Fatal(err)
	}

	if _, err = tx.Exec(`DROP TABLE users;`); err != nil {
		log.Fatal(err)
	}

	if _, err = tx.Exec(`ALTER TABLE users_new RENAME TO users;`); err != nil {
		log.Fatal(err)
	}

	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func seedDefaultUser(db *sql.DB) {
	insertDefaultUser := `
	INSERT INTO users (
		id,
		userName,
		password,
		email,
		city,
		pincode,
		role
	)
	SELECT lower(hex(randomblob(16))), ?, ?, ?, ?, ?, ?
	WHERE NOT EXISTS (
		SELECT 1 FROM users WHERE email = ?
	);`

	defaultEmail := "admin@example.com"
	if _, err := db.Exec(
		insertDefaultUser,
		"admin",
		"admin123",
		defaultEmail,
		"Mumbai",
		"400001",
		"admin",
		defaultEmail,
	); err != nil {
		log.Fatal(err)
	}

	log.Println("\n ========================= \n Default user ensured: \n email: admin@example.com \n Password: admin123 \n ========================= \n")
}
