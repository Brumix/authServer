package utils

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(password string) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(hashedPassword string, testPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err != nil {
		return false
	}
	return true
}
