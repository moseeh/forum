package internal

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err
}
func CompareHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func UUIDGen() string {
	u2, _ := uuid.NewV4()
	return u2.String()
}

func TokenGen(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Token generation failed: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
