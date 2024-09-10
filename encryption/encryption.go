package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

func DeriveEncryptionKey(masterPassword, salt []byte) []byte {
	return pbkdf2.Key(masterPassword, salt, 100000, 32, sha256.New)
}

func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func Encrypt(password []byte, encryptionKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aead.NonceSize())

	encryptedPassword :=aead.Seal(nonce, nonce, password, nil)
	return encryptedPassword, nil
}

func Decrypt(rawEncryptedPassword []byte, encryptionKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aead.NonceSize()
	if len(rawEncryptedPassword) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, encryptedPassword := rawEncryptedPassword[:nonceSize], rawEncryptedPassword[nonceSize:]
	descryptedPassword, err := aead.Open(nil, nonce, encryptedPassword, nil)
	if err != nil {
		return nil, errors.New("could not decipher password")
	}
	return descryptedPassword, nil
}