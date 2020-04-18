package hash

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

func RSAEncodeAndHash(msg, publicKey []byte) (string, error) {
	encoder, err := parsePublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("could not load public key request: %v", err)
	}

	enc, err := encoder.Encode([]byte(msg))
	if err != nil {
		return "", fmt.Errorf("could not encode message: %v", err)
	}
	encodedHash := base64.StdEncoding.EncodeToString(enc)
	encodedHash = strings.Replace(encodedHash, "/", "-", -1)
	encodedHash = strings.Replace(encodedHash, "+", "_", -1)

	return encodedHash, nil
}

func parsePublicKey(pemBytes []byte) (Encoder, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "PUBLIC KEY":
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}

	return newEncoderFromKey(rawkey)
}

func newEncoderFromKey(k interface{}) (Encoder, error) {
	var sshKey Encoder
	switch t := k.(type) {
	case *rsa.PublicKey:
		sshKey = &rsaPublicKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

type Encoder interface {
	Encode(data []byte) ([]byte, error)
}

type rsaPublicKey struct {
	*rsa.PublicKey
}

func (r *rsaPublicKey) Encode(data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, r.PublicKey, data)
}
