package hash

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// SecInDay equals seconds count in day (24 * 3600)
const SecInDay = 86400

// CurrentTimestamp returns md5 hashed current timestamp
func CurrentTimestamp() string {
	hasher := md5.New()

	t, _ := time.Now().MarshalJSON()
	hasher.Write(t)

	return hex.EncodeToString(hasher.Sum(nil))
}

// TokenWithExpTime user to generate token string with expire time
// lifetime - token exp time in seconds
func TokenWithExpTime(lifetime int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Unix() + int64(lifetime),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte("some_secret_string"))

	return tokenString
}

// ValidateToken used to check provided token (expired / correct secrect / etc)
func ValidateToken(token string) bool {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("some_secret_string"), nil
	})

	if err != nil || !t.Valid {
		return false
	}

	return true
}
