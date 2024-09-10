package cmd

import (
	"fmt"
	"strings"

	"github.com/BenjaminPasco/bpass/db"
	"github.com/BenjaminPasco/bpass/encryption"
	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use: "save",
	Short: "Save a given password",
	Long: "Save a given password to an encrypted file",
	Run:  func (cmd *cobra.Command, args []string)  {
		// handling cli inputs
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println(err)
			return
		}
		if password == "" {
			fmt.Println("Undefined password")
			return
		}
		key, err := cmd.Flags().GetString("identifier")
		if err != nil {
			fmt.Println(err)
			return
		}
		if key == "" {
			fmt.Println("Undefined identifier")
			return
		}
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			fmt.Println(err)
			return
		}

		// getting master password
		masterPassword, err := GetMasterPassword()
		if err != nil {
			fmt.Println("Error getting master password: ", err)
			return
		}

		// encrypting password
		salt, err := encryption.GenerateSalt(16)
		if err != nil {
			fmt.Println("Error generating salt")
			return
		}
		encryptionKey := encryption.DeriveEncryptionKey(masterPassword, salt)
		encryptedPassword, err := encryption.Encrypt([]byte(password), encryptionKey)
		if err != nil {
			fmt.Println("Error encrypting password: ", err)
			return
		}

		// saving data to db
		err = db.SavePassword(key, encryptedPassword, salt, description)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "unique") {
				fmt.Println("Error identifier already used")
				return
			}
			fmt.Println("Error saving password:", err)
			return
		}
		fmt.Println("Password saved")
	},
}

func init() {
	saveCmd.Flags().StringP("password", "p", "", "Password to save")
	saveCmd.MarkFlagRequired("password")
	saveCmd.Flags().StringP("identifier", "i", "", "Paremeter to identifie the password")
	saveCmd.MarkFlagRequired("identifier")
	saveCmd.Flags().StringP("description", "d", "", "Optional description like associated email etc...")
}
