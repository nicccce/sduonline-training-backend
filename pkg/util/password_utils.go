package util

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword 用于对密码进行加密
func EncryptPassword(password string) (string, error) {
	if len(password) < 1 {
		return "", errors.New("password must be at least 1 characters")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword 用于验证密码是否正确
func CheckPassword(inputPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}
