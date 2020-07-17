package utils

import (
	"errors"
	"strings"
)

// FormatError formats the error
func FormatError(err string) error {
	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("Password Incorrect")
	}
	return errors.New("Details Incorrect")
}