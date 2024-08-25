package util

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func CheckIdentifier(param string) (resp map[string]interface{}) {

	if strings.Contains(param, "@") {
		resp = map[string]interface{}{"email": param}
	} else {
		resp = map[string]interface{}{"userID": param}
	}
	return resp
}

func HashPassword(password string) (string, error) {
	bytePassword := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)

	return string(hashedPassword), err
}

func CheckPassword(password, inputtedPassword string) bool {
	bytePassword := []byte(password)
	byteInputtedPassword := []byte(inputtedPassword)

	err := bcrypt.CompareHashAndPassword(bytePassword, byteInputtedPassword)

	return err == nil
}
