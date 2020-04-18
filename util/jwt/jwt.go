package jwt

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/djumpen/go-rest-admin/config"
	"github.com/djumpen/go-rest-admin/models"
)

// Error messages
const (
	invalidTokenMsg    = "Invalid token."
	expiredTokenMsg    = "Expired token."
	unparsableTokenMsg = "Unable to parse token."
	unhandableTokenMsg = "Unable to handle token."
)

// GenerateToken used generate new JWT token
func GenerateToken(userID int) string {
	cfg := config.GetJwtConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Unix() + cfg.Lifetime,
		"id":  userID,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(cfg.Secret))

	return tokenString
}

// GetID used to return user id from token claims
func GetID(tokenString string) (models.PKID, error) {
	cfg := config.GetJwtConfig()
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return models.PKID(claims["id"].(float64)), nil
}

// GetString used to return string from token claims
func GetString(tokenString string) (string, error) {
	cfg := config.GetJwtConfig()
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["provider_pms_id"].(string), nil
}

// ValidateToken used for token validation
func ValidateToken(jwtToken string) (bool, error) {
	cfg := config.GetJwtConfig()
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return false, err
	}

	if token.Valid {
		return true, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, errors.New(invalidTokenMsg)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, errors.New(expiredTokenMsg)
		} else {
			return false, errors.New(unparsableTokenMsg)
		}
	}

	return false, errors.New(unhandableTokenMsg)
}
