package cmd

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/BenjaminPasco/bpass/keychain"
)

func GetMasterPassword() ([]byte, error) {
	masterPasswordFromKeychain, err := keychain.RetrieveMasterPassword()
	if err != nil {
		masterPasswordFromInput, err := GenerateMasterPassword()
		if err != nil {
			return nil, err
		}
		err = keychain.StoreMasterPassword(masterPasswordFromInput)
		if err != nil {
			return nil, err
		}
		return masterPasswordFromInput, nil
	}
	return masterPasswordFromKeychain, nil
}

func InputMasterPassword() ([]byte, error) {
	fmt.Print("Enter master password: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return []byte(input), nil
}

func GenerateMasterPassword() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}