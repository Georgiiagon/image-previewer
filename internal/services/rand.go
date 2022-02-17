package services

import (
	"crypto/rand"
	"encoding/base32"
)

func GetToken(length int) string {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

func (s *Service) RandString(length int) string {
	return GetToken(length)
}
