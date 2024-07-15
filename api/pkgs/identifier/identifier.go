package identifier

import (
	"errors"
	"strings"
	"time"
	"unicode"
)

type Identifier[T any] interface {
	New() (T, error)
	FromString(string) (T, error)
	FromBytes([]byte) (T, error)
	Validate(T) (bool, error)
}


type ID string

const (
	alphabet        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits          = "0123456789"
	idLength        = 20
	timestampLength = 12
)

func (i ID) New() ID {
	timestamp := time.Now().Format("060102150405") // YYMMDDhhmmss
	nanoseconds := time.Now().Nanosecond()

	randomPart := make([]byte, idLength-len(string(i))-len(timestamp))
	for idx := range randomPart {
		// Use the nanoseconds value to add some entropy
		if idx%2 == 0 {
			randomPart[idx] = alphabet[(nanoseconds+idx)%len(alphabet)]
		} else {
			randomPart[idx] = digits[(nanoseconds+idx)%len(digits)]
		}
	}

	id := string(i) + timestamp + string(randomPart)
	return ID(id)
}

func (i ID) FromString(s string) (ID, error) {
	id := ID(s)
	if !id.Validate() {
		return "", errors.New("invalid identifier format")
	}
	return id, nil
}

func (i ID) FromBytes(b []byte) (ID, error) {
	id := ID(b)
	if !id.Validate() {
		return "", errors.New("invalid identifier format")
	}
	return id, nil
}

func (i ID) Validate() bool {
	idStr := string(i)

	// Ensure the length is correct
	if len(idStr) != idLength {
		return false
	}

	// Extract the prefix, timestamp, and random part
	// Find the first underscore to determine the prefix
	underscoreIndex := strings.Index(idStr, "_")
	if underscoreIndex == -1 {
		return false
	}
	remaining := idStr[underscoreIndex+1:]

	// Ensure the remaining part length is at least as long as the timestamp
	if len(remaining) < timestampLength {
		return false
	}

	// Extract the timestamp and random part from the remaining part
	timestamp := remaining[:timestampLength]
	randomPart := remaining[timestampLength:]

	// Validate the timestamp part
	if _, err := time.Parse("060102150405", timestamp); err != nil {
		return false
	}

	// Validate the random part: alternating letters and digits
	for idx, char := range randomPart {
		if idx%2 == 0 {
			if !unicode.IsLetter(char) || !contains(alphabet, char) {
				return false
			}
		} else {
			if !unicode.IsDigit(char) || !contains(digits, char) {
				return false
			}
		}
	}

	// If all checks pass, the ID is valid
	return true
}

// contains checks if a rune is in a string
func contains(s string, r rune) bool {
	for _, char := range s {
		if char == r {
			return true
		}
	}
	return false
}
