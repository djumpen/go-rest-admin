package hash

import (
	"testing"
)

func TestValidPassword(t *testing.T) {
	password := "123123123"

	generated, err := GeneratePassword(password)

	if err != nil {
		t.Error(err.Error())
	}

	if !ValidPassword(password, generated) {
		t.Error("Incorrect password generation / validation")
	}
}
