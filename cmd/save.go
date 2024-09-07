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

	"github.com/BenjaminPasco/bpass/db"
	"github.com/spf13/cobra"
)

var encryptionKey = []byte("abababababababababababababababab")

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
			fmt.Println("Error encrypting password: ", err)
			return
		}
		err = savePassword(key, encryptedPassword)
		if err != nil {
			fmt.Println("Error saving password:", err)
			return
		}
		fmt.Println("Password saved")
	},
}

func init() {
	saveCmd.Flags().StringP("password", "p", "", "Password to save")
	saveCmd.MarkFlagRequired("password")
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