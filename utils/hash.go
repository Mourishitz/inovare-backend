package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckValidHash(pw string) bool {
	return len(pw) > 0 && len(pw) < 100 && (pw[:4] == "$2a$" || pw[:4] == "$2b$" || pw[:4] == "$2y$")
}

func CheckValidHashWithPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
