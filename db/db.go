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

type PasswordData struct {
	Identifier string
	Description string
}

func ListPasswords() ([]PasswordData, error) {
	listPasswordsStatement := `
		SELECT key, description FROM passwords
	`
	rows, err := DB.Query(listPasswordsStatement)
	if err != nil {
		return nil, err
	}

	// processing raw data
	var entries []PasswordData
	for rows.Next() {
		var entry PasswordData
		err := rows.Scan(&entry.Identifier, &entry.Description)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func DeleteAllPasswords() error {
	
	deleteAllPasswordsStatement := `
		DELETE FROM passwords;
	`

	_, err := DB.Exec(deleteAllPasswordsStatement)
	return err
}

func DeletePassword(identifier string) error {
	
	deletePasswordStatement := `
		DELETE FROM passwords WHERE key = ?;
	`

	result, err := DB.Exec(deletePasswordStatement, identifier)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows ==  0 {
		return errors.New("could not find password with matching identifier")
	}
	return err
}