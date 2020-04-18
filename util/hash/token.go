package hash

import (
	"crypto/rand"
	"fmt"
)

// GenerateToken generates random 7 character token.
func GenerateToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
