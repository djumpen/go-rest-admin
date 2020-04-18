package validation

import (
	"regexp"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

// checkDate check date "YYYY-mm-dd" format
func checkDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	if date == "" {
		return true
	}
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// checkDate check name to contain only
func checkName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	isCorrect := regexp.MustCompile(`^[A-Za-z-— ]+$`).MatchString
	return isCorrect(name)
}
