package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// Bcrypt generate a hash based on the given string
func Bcrypt(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

// CompareHash returns true if the hash and password match
func CompareHash(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
