package signature

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

var (
	AesKey = []byte("V6-F6YChZ-esxVybB22UKeyawU6XYDyK")
)

// AesEncrypt aes encrypt
func AesEncrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(AesKey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(plaintext))
	return hex.EncodeToString(ciphertext), nil
}

// AesDecrypt aes decrypt
func AesDecrypt(d string) (string, error) {
	ciphertext, err := hex.DecodeString(d)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(AesKey)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
