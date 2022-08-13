package utils

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "username") {
		return errors.New("Kullanıcı adı alınmış")
	}
	if strings.Contains(err, "email") {
		return errors.New("Email Adresi alınmış")
	}
	return nil
}
