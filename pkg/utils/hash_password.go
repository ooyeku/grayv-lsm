package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a password as input and returns the hashed password as a string.
// It uses the bcrypt algorithm to generate a secure hash from the input password.
// The function returns the hashed password string and any error that occurred during the hashing process.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a given password with its hashed counterpart and
// returns true if they match; otherwise, it returns false.
// It uses bcrypt.CompareHashAndPassword to compare the hashed password with the
// plain-text password.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
