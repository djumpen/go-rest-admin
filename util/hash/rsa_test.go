package hash

import (
	"testing"
)

func TestRSADecrypt(t *testing.T) {
	text := "CijhoapFevITGUX2cO0THOuxTB-Bql2bRwFrP7bjK3zdZIOCr-jyLhEZIYgiZyKzte38PWS-MaNJOEqGUY161w=="
	key := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAJbnNVqc0qWayz2jCDesrh64JJo1BuHlwOQppcu2QaIb/4r7+wrz
sNUgRm2GY1l+RNoUIWqf5G25/3/Hh4I+91cCAwEAAQJASfQpq6yrd0lzTVO21UIl
Wxy3o3NVWiP08lyOylUZuU3mQA5aK+rqdbaE7G7llC8RSgVGVD1HbmVxgLE7P2Qg
EQIhAPO+C779SNcIp5NMVPqoQTlgldSkfsIfz0VQQUTDXevTAiEAnn3x5AIFiY/C
cE/7Xj4Cm8gge6kzAbkdV2msEuKXp+0CIQDLoZpulWylObXGeZ8FSkwzg12pqUO9
KpYfck0VBaMRwQIgDKj80HzE2ncsTfJlnuKPLMhwp9AdLe8Og/QB9cQ53wUCIHg+
2I14N0H5UdYW6u4soLvqyK/643m3M81fzyh34MbO
-----END RSA PRIVATE KEY-----
`)

	k, err := RSADecrypt(text, key)
	if err != nil {
		t.Errorf("Unexpected error - %s", err.Error())
	}

	if string(k) != "3588c54e-88ef-4700-94fe-7f45c1b71329" {
		t.Error("Wrong value returned")
	}

	if _, err = RSADecrypt(text, []byte("wrongkey")); err == nil {
		t.Error("Key error expected, got nil")
	}
}
