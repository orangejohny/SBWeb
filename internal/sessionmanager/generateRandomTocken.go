package sessionmanager

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomBytes returns slice of random bytes
func GenerateRandomBytes(length int) ([]byte, error) {
	tocken := make([]byte, length)
	_, err := rand.Read(tocken)
	if err != nil {
		return nil, err
	}

	return tocken, nil
}

// GenerateRandomString returns random string that can be used as a tocken
func GenerateRandomString(length int) (string, error) {
	tocken, err := GenerateRandomBytes(length)
	return base64.URLEncoding.EncodeToString(tocken), err
}
