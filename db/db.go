package db

import (
	"database/sql"
	"errors"
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
		encrypted_password BYTEA NOT NULL,
		salt BYTEA NOT NULL,
		description TEXT
	);`

	_, err = DB.Exec(createPasswordTableStatement)
	if err != nil {
		fmt.Println("Error setting up passwords table", err)
	}
	return err
}

func SavePassword(key string, encryptedPassword []byte, salt []byte, description string) error {
	
	savePasswordStatement := `
		INSERT INTO passwords (
			key,
			encrypted_password,
			salt,
			description
		)
		VALUES (?, ?, ?, ?)
	`
	_, err := DB.Exec(savePasswordStatement, key, encryptedPassword, salt, description)
	return err
}

func GetEnryptedPasswordFromDb(key string) ([]byte, []byte, error) {
	
	getPasswordStatement := `
		SELECT encrypted_password, salt FROM passwords WHERE key = ?
	`

	result, err := DB.Query(getPasswordStatement, key)
	if err != nil {
		return nil, nil, err
	}
	if !result.Next() {
		return nil, nil, errors.New("no password match this identifier")
	}
	var encryptedPassword []byte
	var salt []byte
	result.Scan(&encryptedPassword, &salt)
	return encryptedPassword, salt, nil
}