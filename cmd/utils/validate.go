package utils

import "regexp"

func ValidateEmail(email string) bool {
	if len(email) > MaxIdentityLen || len(email) < MinIdentityLen {
		return false
	}
	match, _ := regexp.MatchString(EmailRgx, email)
	return match
}

func ValidatePassword(password string) bool {
	if len(password) > MaxIdentityLen || len(password) < MinIdentityLen {
		return false
	}
	return true
}

func ValidateName(name string) bool {
	if len(name) > MaxNameLen || len(name) < MinNameLen {
		return false
	}
	return true
}
