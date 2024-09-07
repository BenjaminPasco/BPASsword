package cmd

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var passwordFile = "password.txt"
var encryptionKey = []byte("abababababababababababababababab")

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
			fmt.Println("Empty password")
			return
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("account/application/domain: ")
		key, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		key = key[:len(key)-1]

		encryptedPassword, err := encrypt(password)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = SavePassword(key, encryptedPassword)
		if err != nil {
			fmt.Println("Error saving password:", err)
			return
		}
		fmt.Println("Password saved")
	},
}

func init() {
	saveCmd.Flags().StringP("password", "p", "", "Password to save")
}

func encrypt(password string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
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

    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func SavePassword(key string, encryptedPassword string) error {
	file, err := os.OpenFile(passwordFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s: %s\n", key, encryptedPassword))
	return err
}