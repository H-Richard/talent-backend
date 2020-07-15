package utils

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "username") {
		return errors.New("Username Already Taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("Username Already Taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("Password Incorrect")
	}
	return errors.New("Details Incorrect")
}