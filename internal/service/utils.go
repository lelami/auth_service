package service

import (
	"crypto/rand"
	"encoding/base64"
	mr "math/rand"
	"strconv"
	"time"
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

func randOTP() string {
	r := mr.New(mr.NewSource(time.Now().UnixNano()))
	randomNumber := r.Intn(9000) + 1000
	rd := strconv.Itoa(randomNumber)

	return rd
}
