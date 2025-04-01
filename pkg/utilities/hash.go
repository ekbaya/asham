package utilities

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plain-text password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	charsetLength := len(charset)

	// Create a byte slice for the password
	password := make([]byte, length)

	// Fill it with random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Map random bytes to characters in our charset
	for i := 0; i < length; i++ {
		password[i] = charset[int(randomBytes[i])%charsetLength]
	}

	return string(password), nil
}
