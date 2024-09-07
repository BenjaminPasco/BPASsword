package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		fmt.Println("Error opening db", err)
	}

	createPasswordTableStatement := `
	CREATE TABLE IF NOT EXISTS passwords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT NOT NULL UNIQUE,
		encrypted_password TEXT NOT NULL
	);`

	_, err = DB.Exec(createPasswordTableStatement)
	if err != nil {
		fmt.Println("Error setting up passwords table", err)
	}
	return err
}