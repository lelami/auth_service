package service

import (
	"crypto/rand"
	"encoding/base64"
)

const passwordLength = 15

func randPsw() (string, error) {
	byteLength := (passwordLength*6 + 7) / 8

	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	password := base64.RawURLEncoding.EncodeToString(bytes)

	if len(password) > passwordLength {
		password = password[:passwordLength]
	}

	if len(password) < passwordLength {
		password += "0"
	}

	return password, nil
}
