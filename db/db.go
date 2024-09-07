package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./passwords.db")
	if err != nil {
		fmt.Println("Error opnening db", err)
	}

	createPasswordTableStatement := `
	CREATE TABLE IF NOT EXISTS passwords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT NOT NULL,
		encrypted_password TEXT NOT NULL
	);`

	_, err = DB.Exec(createPasswordTableStatement)
	if err != nil {
		fmt.Println("Error setting up passwords table", err)
	}
	return err
}