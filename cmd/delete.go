package cmd

import (
	"fmt"

	"github.com/BenjaminPasco/bpass/db"
	"github.com/spf13/cobra"
)

var deleteAll bool

var deleteCmd = &cobra.Command{
	Use: "delete",
	Short: "delete passwords",
	Run:  func (cmd *cobra.Command, args []string)  {
		identifier, err := cmd.Flags().GetString("identifier")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if len(identifier) > 0 && deleteAll {
			fmt.Println("Error: cannot use --identifier and --all at the same time")
			return 
		}
		if deleteAll {
			err := db.DeleteAllPasswords()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			fmt.Println("Successfully deleted all passwords")
			return
		}
		err = db.DeletePassword(identifier)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Successfully deleted: %s\n", identifier)
	},
}

func init() {
	deleteCmd.Flags().StringP("identifier", "i", "", "Paremeter to identifie the password, mutually exclusive with the --all flag")
	deleteCmd.Flags().BoolVarP(&deleteAll, "all", "a", false, "Parameter to delete all passwords at once, mutually exclusive with --identifier")
}
