package utils

import (
	"regexp"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/constants"
)

func ValidateEmail(email string) bool {
	if len(email) > constants.MaxIdentityLen || len(email) < constants.MinIdentityLen {
		return false
	}
	match, _ := regexp.MatchString(constants.EmailRgx, email)
	return match
}

func ValidatePassword(password string) bool {
	if len(password) > constants.MaxIdentityLen || len(password) < constants.MinIdentityLen {
		return false
	}
	return true
}

func ValidateName(name string) bool {
	if len(name) > constants.MaxNameLen || len(name) < constants.MinNameLen {
		return false
	}
	return true
}
