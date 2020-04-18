package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GeneratePassword is a wrapper around default Goland bcrypt hash library
func GeneratePassword(password string) (string, error) {
	var hash []byte
	var err error

	if hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return "", err
	}

	return string(hash), nil
}

// ValidPassword used to validate user password (compare string/hashed) using default Golang bcrypt library
func ValidPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println(err)
	return err == nil
}
