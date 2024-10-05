package utils

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(strings.ToLower(email))
}

func ValidateURL(url string) bool {
	// This is a simple check. For production, consider using a more robust solution.
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func ValidateNonEmptyString(s string) bool {
	return strings.TrimSpace(s) != ""
}

func ValidateMaxLength(s string, maxLength int) bool {
	return len(s) <= maxLength
}
