package idgen

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"strings"
	"time"
)

type Identifier[T any] interface {
	// New creates a new random Identifier
	New() (T, error)
	// FromString parses a string representation of the Identifier
	FromString(string) (T, error)
	// FromBytes parses a byte array representation of the Identifier
	FromBytes([]byte) (T, error)
	// Validate the Identifier conforms to the specification and is valid
	Validate(T) (bool, error)
}

type CustomID string

const (
	idLength     = 20
	timeBytes    = 6
	randomBytes  = 9
	encodedChars = "0123456789ABCDEFGHJKMNPQRSTVWXYZ" // 32 characters, excluding I, L, O, U
)

var encoding = base32.NewEncoding(encodedChars).WithPadding(base32.NoPadding)

// CustomIDHandler implements the Identifier interface for CustomID
type CustomIDHandler struct{}

func (h CustomIDHandler) New() (CustomID, error) {
	// Get current timestamp (48 bits, millisecond precision)
	now := time.Now().UnixNano() / 1e6
	timeComponent := make([]byte, timeBytes)
	for i := timeBytes - 1; i >= 0; i-- {
		timeComponent[i] = byte(now & 0xFF)
		now >>= 8
	}

	// Generate random component
	randomComponent := make([]byte, randomBytes)
	_, err := rand.Read(randomComponent)
	if err != nil {
		return "", err
	}

	// Combine time and random components
	combined := append(timeComponent, randomComponent...)

	// Encode to base32
	encoded := encoding.EncodeToString(combined)

	return CustomID(encoded), nil
}

func (h CustomIDHandler) FromString(s string) (CustomID, error) {
	if len(s) != idLength {
		return "", errors.New("invalid identifier length")
	}
	return CustomID(s), nil
}

func (h CustomIDHandler) FromBytes(b []byte) (CustomID, error) {
	if len(b) != idLength {
		return "", errors.New("invalid identifier length")
	}
	return CustomID(string(b)), nil
}

func (h CustomIDHandler) Validate(id CustomID) (bool, error) {
	if len(id) != idLength {
		return false, errors.New("invalid identifier length")
	}

	// Check if all characters are valid
	for _, char := range id {
		if !strings.ContainsRune(encodedChars, char) {
			return false, errors.New("invalid character in identifier")
		}
	}

	return true, nil
}

//func (id CustomID) String() string {
//	return string(id)
//}
//
//func (id CustomID) Bytes() []byte {
//	return []byte(id)
//}
//
//// Ensure CustomIDHandler implements Identifier interface
//var _ Identifier[CustomID] = (*CustomIDHandler)(nil)
