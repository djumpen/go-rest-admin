package auth

import (
	"errors"

	"github.com/djumpen/go-rest-admin/apperrors"
	"github.com/djumpen/go-rest-admin/config"
	"github.com/djumpen/go-rest-admin/util/hash"
	"github.com/djumpen/go-rest-admin/util/jwt"

	"strings"

	"github.com/gin-gonic/gin"
)

//API messages set
const (
	msgTokenRequired    = "missed bearer header"
	msgNotIntegerValue  = "value should be an integer"
	msgTokenCheckFailed = "invalid token"
)

func RequireJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if len(authHeader) <= 0 {
			c.Error(apperrors.NewUnauthorized(errors.New(msgTokenRequired)))
		}
		token := strings.Split(authHeader, " ")
		if len(token) < 2 {
			c.Error(apperrors.NewUnauthorized(errors.New(msgTokenCheckFailed)))
		}
		if success, err := jwt.ValidateToken(token[1]); !success {
			c.Error(apperrors.NewUnauthorized(err))
		}
	}
}

func decodeCID(CID string) (string, error) {
	id, err := hash.RSADecrypt(CID, config.RSAPrivateKey())
	if err != nil {
		return "", apperrors.NewUnauthorized(errors.New(msgTokenCheckFailed))
	}
	return string(id), nil
}
