package cmd

import (
	"fmt"

	"github.com/BenjaminPasco/bpass/db"
	"github.com/BenjaminPasco/bpass/encryption"
	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use: "decode",
	Short: "return decoded password",
	Long: "return the decoded password for a given identifier",
	Run:  func (cmd *cobra.Command, args []string)  {
		// handling cli inputs
		key, err := cmd.Flags().GetString("identifier")
		if err != nil {
			fmt.Println(err)
			return
		}
		if key == "" {
			fmt.Println("Undefined identifier")
			return
		}

		// getting master password
		masterPassword, err := GetMasterPassword()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// getting data from db
		encryptedPassword, salt, err := db.GetEnryptedPasswordFromDb(key)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// generation encryption key
		encryptionKey :=encryption.DeriveEncryptionKey(masterPassword, salt)

		// deciphering password
		decryptedPassword, err := encryption.Decrypt(encryptedPassword, encryptionKey)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("password: ", string(decryptedPassword))
	},
}

func init(){
	decodeCmd.Flags().StringP("identifier", "i", "", "Paremeter to identifie the password")
	decodeCmd.MarkFlagRequired("identifier")
}
