/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/BenjaminPasco/bpass/cmd"
	"github.com/BenjaminPasco/bpass/db"
)

func main() {
	err := db.InitDB() // start DB connection, create table if needed
	if err != nil {
		fmt.Println("main", err)
	}
	defer db.DB.Close() // defer DB connection close
	cmd.Execute()
}
