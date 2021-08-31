package crypto

import (
	"crypto/sha256"
) //NewSHA256 method for hash

func NewSHA256() []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func CheckHash(data []byte, hash string) bool {
	hash := string(NewSHA256())
	uuid := data.uuid

	return false
}
