package cmd

import (
	"fmt"
	"strings"

	"github.com/BenjaminPasco/bpass/db"
	"github.com/spf13/cobra"
)

var list = &cobra.Command{
	Use: "list",
	Short: "list passwords",
	Run:  func (cmd *cobra.Command, args []string)  {
		entries, err := db.ListPasswords()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		const format = "%-20s %-20s %-20s\n"
		fmt.Printf(format, "Identifier", "Password", "Description")
		fmt.Println(strings.Repeat("-", 60))
		for _, entry := range entries {
			fmt.Printf(format, entry.Identifier, strings.Repeat("*", 18), entry.Description)
		}
	},
}
