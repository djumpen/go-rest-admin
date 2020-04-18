package validation

import (
	"errors"
	"regexp"
)

// StripValidatePhone allows all symbols except numbers and check lenght
func StripValidatePhone(p string) (string, error) {
	const phoneLen = 10
	var re = regexp.MustCompile(`[^0-9]*`)
	p = re.ReplaceAllString(p, "")
	if p[:1] == "1" {
		p = p[1:]
	}
	if len(p) < phoneLen {
		return "", errors.New("Phone is not valid")
	}
	return p, nil
}
