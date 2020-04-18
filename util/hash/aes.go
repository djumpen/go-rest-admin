package hash

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

// Encrypt used to encrypt text by provided key/iv via AES algorithm with 128-bit key
func Encrypt(keyStr, ivStr, text string) (string, error) {
	key := []byte(keyStr)
	iv := []byte(ivStr)

	plaintext := pad([]byte(text))

	if len(plaintext)%aes.BlockSize != 0 {
		return "", errors.New("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return base64Encode(ciphertext), nil
}

// Decrypt used to decrypt text by provided key/iv via AES algorithm with 128-bit key
func Decrypt(keyStr, ivStr, text string) (string, error) {
	key := []byte(keyStr)
	iv := []byte(ivStr)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decoded, err := base64Decode(text)
	if err != nil {
		return "", err
	}

	ciphertext := decoded[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	unpaded, err := unpad(ciphertext)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", unpaded), nil
}

// base64Encode input byte array string
func base64Encode(toEncode []byte) string {
	return base64.StdEncoding.EncodeToString(toEncode)
}

// base64Decode input string
func base64Decode(toDecode string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(toDecode)
	return decoded, err
}

// pad used to add more symbols to ciphertext if len(ciphertext)%aes.BlockSize
func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(src, padtext...)
}

// unpad used to remove additional symbols from ciphertext (added if len(ciphertext)%aes.BlockSize)
func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}
