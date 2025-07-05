package passwordhashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPasswordArgon2(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	encoded := fmt.Sprintf("%s$%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash))
	return encoded, nil
}

func ComparePasswordArgon2(inputPassword, storedHash string) (bool, error) {
	parts := strings.Split(storedHash, "$")
	if len(parts) != 2 {
		return false, errors.New("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	inputHash := argon2.IDKey([]byte(inputPassword), salt, 1, 64*1024, 4, 32)

	// Constant time compare to avoid timing attacks
	if subtle.ConstantTimeCompare(expectedHash, inputHash) == 1 {
		return true, nil
	}
	return false, nil
}
