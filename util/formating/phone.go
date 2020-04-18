package formating

import (
	"strings"
)

// FormatPhoneToEmail formats 5555555555 phone number to 555-555-5555 format or returns given if phone less than 10 digits
func FormatPhoneToEmail(phone string) string {
	if len(phone) < 10 {
		return phone
	}
	return strings.Join([]string{phone[:3], phone[3:6], phone[6:]}, "-")
}
