// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package sessionmanager

import (
	"crypto/rand"
	"encoding/base64"
)

// generateRandomBytes returns slice of random bytes
func generateRandomBytes(length int) ([]byte, error) {
	tocken := make([]byte, length)
	_, err := rand.Read(tocken)
	if err != nil {
		return nil, err
	}

	return tocken, nil
}

// generateRandomString returns random string that can be used as a tocken
func generateRandomString(length int) (string, error) {
	tocken, err := generateRandomBytes(length)
	return base64.URLEncoding.EncodeToString(tocken), err
}
