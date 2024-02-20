package helper

import (
	"golang.org/x/crypto/bcrypt"
)

var hasherFunc = bcrypt.GenerateFromPassword

func HashPassword(password string) (string, error) {
	hash, err := hasherFunc([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
