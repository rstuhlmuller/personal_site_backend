package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func Truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}
