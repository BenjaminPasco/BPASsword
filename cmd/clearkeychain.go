package cmd

import (
	"fmt"

	"github.com/BenjaminPasco/bpass/keychain"
	"github.com/spf13/cobra"
)

var clearKeychain = &cobra.Command{
	Use: "clearkeychain",
	Short: "delete master key in keychain",
	Run:  func (cmd *cobra.Command, args []string)  {
		err := keychain.DeleteMasterPassword()
		if err != nil {
			fmt.Println("Error deleting master key from keychain: ", err)
			return 
		}
		fmt.Println("Successfully cleared master key from keychain")
	},
}
