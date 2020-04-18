package hash

import (
	"testing"
	"time"
)

func TestCurrentTimestamp(t *testing.T) {
	h1 := CurrentTimestamp()
	time.Sleep(time.Second) //GenerateFilename will hash timestamp, let's wait 1sec to be sure
	h2 := CurrentTimestamp()

	if h1 == h2 {
		t.Errorf("hash1 shouldn't equals hash2, smth goes wrong")
	}
}

func TestTokenWithExpTime(t *testing.T) {
	t1 := TokenWithExpTime(1)
	t2 := TokenWithExpTime(2)

	if t1 == t2 {
		t.Error("Token shouldn't be the same")
	}
}

func TestValidateToken(t *testing.T) {
	t1 := TokenWithExpTime(1)
	time.Sleep(time.Second * 2)

	if ValidateToken(t1) {
		t.Error("Token is valid, but should be expired")
	}

	t2 := TokenWithExpTime(1000)

	if !ValidateToken(t2) {
		t.Errorf("Token is invalid, but should be")
	}
}
