/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BenjaminPasco/bpass/cmd"
	"github.com/BenjaminPasco/bpass/db"
)

func getDBFilePath() (string, error) {
	switch system := os.Getenv("OS"); system {
	case "Windows_NT":
		return filepath.Join(os.Getenv("APPDATA"), "bpass", "password.db"), nil
	default:
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, "Library", "Application Support", "bpass", "password.db"), nil
	}
}

func ensureDirExists(dir string) error {
	_, err := os.Stat(dir)
	
	if errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	dbPath, err := getDBFilePath() // get the path of the file that will store passwords
	if err != nil {
		fmt.Println("Error: could not get db file location", err)
		return
	}
	dir := filepath.Dir(dbPath)
	err = ensureDirExists(dir)
	if err != nil {
		fmt.Println("Error: directory for db file does not exist or could not be created", err)
		return
	}
	err = db.InitDB(dbPath) // start DB connection, create table if needed
	if err != nil {
		fmt.Println("Error: could not open connection to db file", err)
		return
	}
	defer db.DB.Close() // defer DB connection close
	cmd.Execute()
}
