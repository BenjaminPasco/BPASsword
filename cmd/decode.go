package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/BenjaminPasco/bpass/constants"
	"github.com/BenjaminPasco/bpass/db"
	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use: "decode",
	Short: "return decoded password",
	Long: "return the decoded password for a given identifier",
	Run:  func (cmd *cobra.Command, args []string)  {
		key, err := cmd.Flags().GetString("identifier")
		if err != nil {
			fmt.Println(err)
			return
		}
		if key == "" {
			fmt.Println("Undefined identifier")
			return
		}
		encodedPassword, err := getEncodedPasswordFromDb(key)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		decodedPassword, err := decodePassword([]byte(encodedPassword))
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println("password: ", decodedPassword)
	},
}

func init(){
	decodeCmd.Flags().StringP("identifier", "i", "", "Paremeter to identifie the password")
	decodeCmd.MarkFlagRequired("identifier")
}

func getEncodedPasswordFromDb(key string) (string, error) {
	
	getPasswordStatement := `
		SELECT encrypted_password FROM passwords WHERE key = ?
	`

	result, err := db.DB.Query(getPasswordStatement, key)
	if err != nil {
		return "", err
	}
	if !result.Next() {
		return "", errors.New("no password match this identifier")
	}
	var dest string
	result.Scan(&dest)
	return dest, nil
}

func decodePassword(base64EncodedPassword []byte) (string, error){

    encodedPassword, err := base64.StdEncoding.DecodeString(string(base64EncodedPassword))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(constants.EncryptionKey))
	if err != nil {
		return "", err
	}
	if len(encodedPassword) < aes.BlockSize {
		return "", errors.New("encrypted block size if too short")
	}

	iv := encodedPassword[:aes.BlockSize]
	decodedPassword := encodedPassword[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(decodedPassword, decodedPassword)
	return string(decodedPassword), nil
}