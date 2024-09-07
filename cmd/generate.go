package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/spf13/cobra"
)

// available flags
var length int

// this is the command
var generateCmd = &cobra.Command{
	Use: "generate",
	Short: "Generate a random password",
	Long: "Generate a random password of variable length, 24 being the default, the password will contain uppercases, lowercases, digits and numbers",
	Run:  func (cmd *cobra.Command, args []string)  {
		password, err := GeneratePassword(length)
		if err != nil {
			fmt.Println("Error generating password:", err)
			return
		}
		fmt.Println("Generated Password:", password)
	},
}

func init() {
	// rootCmd.AddCommand(generateCmd)

	// define the length flag for the command with the default value of 24
	generateCmd.Flags().IntVarP(&length, "length", "l", 24, "Length of the password")
}

// this is the function generating the password
func GeneratePassword(length int) (string, error) {
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	special := "!@#$%^&*()-_=+{}[]|;:,.<>?/~`"

	charSet := lower + upper + digits + special

	var password []byte
	validPassword := false
	// we check if there is a symbol in the generated password
	for !validPassword {
		for i := 0; i <length; i++ {
			char, err := randomChar(charSet)
			if err != nil {
				return "", err
			}
			if strings.Contains(special, string(char)){
				validPassword = true
			}
			password = append(password, char)
		}
	}
	return string(password), nil
}

func randomChar(charSet string) (byte, error) {
	max := big.NewInt(int64(len(charSet)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return charSet[n.Int64()], nil

}