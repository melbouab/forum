// package helpers - VerifyPassword (corrected, minor improvements for clarity)
package helpers

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/argon2"
)

func VerifyPassword(password, encodedHash string) error {
	parts := strings.Split(encodedHash, ".")
	if len(parts) != 2 {
		return errors.New("invalid encoded hash")
	}

	saltBase64 := parts[0]
	hashPasswordBase64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return errors.New("failed to decode the salt")
	}

	hashPassword, err := base64.StdEncoding.DecodeString(hashPasswordBase64)
	if err != nil {
		return errors.New("failed to decode the hash")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	if len(hashPassword) != len(hash) {
		return errors.New("incorrect password")
	}

	if subtle.ConstantTimeCompare(hash, hashPassword) == 1 {
		return nil
	}
	return errors.New("incorrect password")
}
