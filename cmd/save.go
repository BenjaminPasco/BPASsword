package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/BenjaminPasco/bpass/constants"
	"github.com/BenjaminPasco/bpass/db"
	"github.com/spf13/cobra"
)

type PasswordEntry struct {
	Key string `json:"key"`
	Password string `json:"password"`
}

var saveCmd = &cobra.Command{
	Use: "save",
	Short: "Save a given password",
	Long: "Save a given password to an encrypted file",
	Run:  func (cmd *cobra.Command, args []string)  {
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

		encryptedPassword, err := encrypt(password)
		if err != nil {
			fmt.Println("Error encrypting password: ", err)
			return
		}
		err = savePassword(key, encryptedPassword)
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
}

func encrypt(password string) (string, error) {
	block, err := aes.NewCipher([]byte(constants.EncryptionKey))
	if err != nil {
		return "", err
	}

	passwordBytes := []byte(password)
    ciphertext := make([]byte, aes.BlockSize+len(passwordBytes))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], passwordBytes)

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func savePassword(key string, encryptedPassword string) error {
	
	savePasswordStatement := `
		INSERT INTO passwords (
			key,
			encrypted_password
		)
		VALUES (?, ?)
	`
	_, err := db.DB.Exec(savePasswordStatement, key, encryptedPassword)
	return err
}