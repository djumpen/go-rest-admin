package hash

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	text := "b1nkq9b3oht0c255v4r0"
	key := "LKHlhb899Y09olUi"
	iv := "MMMMMMMMMMMMMMMM"

	encrypted, err := Encrypt(key, iv, text)
	if err != nil {
		t.Error(err.Error())
	}

	expected := "AAAAAAAAAAAAAAAAAAAAADUF9gn+7+b9ZRZQxoJ7Hu+bi1DjMwcZN0Jkna9cl9Bq"

	if encrypted != expected {
		t.Error("Invalid encryption")
	}
}

func TestDecrypt(t *testing.T) {
	text := "AAAAAAAAAAAAAAAAAAAAADUF9gn+7+b9ZRZQxoJ7Hu+bi1DjMwcZN0Jkna9cl9Bq"
	key := "LKHlhb899Y09olUi"
	iv := "MMMMMMMMMMMMMMMM"

	decrypted, err := Decrypt(key, iv, text)
	if err != nil {
		t.Error(err.Error())
	}

	expected := "b1nkq9b3oht0c255v4r0"

	if decrypted != expected {
		t.Error("Invalid decryption")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	text := "b1nkq9b3oht0c255v4r0"
	key := "LKHlhb899Y09olUi"
	iv := "MMMMMMMMMMMMMMMM"

	encrypted, err := Encrypt(key, iv, text)
	if err != nil {
		t.Error(err.Error())
	}

	decrypted, err := Decrypt(key, iv, encrypted)
	if err != nil {
		t.Error(err.Error())
	}

	if decrypted != text {
		t.Error("Invalid decryption of encrypted text")
	}
}
