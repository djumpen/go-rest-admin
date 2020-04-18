package hash

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

// RSADecrypt is used to decode encoded string using RSA private key
func RSADecrypt(text string, key []byte) ([]byte, error) {
	b1, err := base64Dec(text)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, b1)
}

// base64Dec is used to convert base64-encoded value to string
func base64Dec(s1 string) ([]byte, error) {
	s1 = strings.Replace(s1, "_", "+", -1)
	s1 = strings.Replace(s1, " ", "+", -1)
	s1 = strings.Replace(s1, "-", "/", -1)
	s1 = strings.Replace(s1, "\n", "", -1)
	s1 = strings.Replace(s1, "\r", "", -1)
	return base64.StdEncoding.DecodeString(s1)
}
