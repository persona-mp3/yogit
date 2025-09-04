package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
)

// "fmt"
// "v1/common"

func Hasher(content []byte) string {
	hasher := sha1.New()
	if _, err := hasher.Write(content); err != nil {
		log.Fatalf("error: occured in hashing content\n %s\n", err)
	}
	hashedContent := hasher.Sum(nil)

	hashId := hex.EncodeToString(hashedContent)
	return hashId
}

