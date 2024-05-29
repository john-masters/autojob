package utils

import "regexp"

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	lowerCase := regexp.MustCompile(`[a-z]`)
	upperCase := regexp.MustCompile(`[A-Z]`)
	numeric := regexp.MustCompile(`\d`)
	specialChar := regexp.MustCompile(`[@$!%*?&]`)

	if !lowerCase.MatchString(password) {
		return false
	}

	if !upperCase.MatchString(password) {
		return false
	}

	if !numeric.MatchString(password) {
		return false
	}

	if !specialChar.MatchString(password) {
		return false
	}

	return true
}
