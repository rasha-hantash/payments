package identifier

import (
	// "math/rand"
	"errors"
	"math/big"
	"strings"
	"time"
)

const (
	timestampLength = 10
	randomLength    = 9
	separatorChar   = '-'
	base            = 62
)

var (
	ErrInvalidFormat = errors.New("invalid identifier format")
	ErrInvalidLength = errors.New("invalid identifier length")
)

type Identifier[T any] interface {
	New() (T, error)
	FromString(string) (T, error)
	FromBytes([]byte) (T, error)
	Validate(T) (bool, error)
}

// todo turn this into a struct that stores the UUID, NEW() will convert to human readable string
type ID string

var charset = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits   = "0123456789"
	idLength = 20
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
	if ok, err := id.Validate(id); !ok {
		return "", err
	}
	return id, nil
}

func (i ID) FromBytes(b []byte) (ID, error) {
	return ID(b), nil
}

func (i ID) Validate(id ID) (bool, error) {
	if len(id) != timestampLength+1+randomLength {
		return false, ErrInvalidLength
	}

	parts := strings.Split(string(id), string(separatorChar))
	if len(parts) != 2 {
		return false, ErrInvalidFormat
	}

	timestamp, random := parts[0], parts[1]

	if len(timestamp) != timestampLength || len(random) != randomLength {
		return false, ErrInvalidLength
	}

	_, err := decodeBase62(timestamp)
	if err != nil {
		return false, err
	}

	_, err = decodeBase62(random)
	if err != nil {
		return false, err
	}

	return true, nil
}

func encodeBase62(num *big.Int) string {
	encoded := make([]byte, 0, 10)
	for num.Sign() > 0 {
		mod := new(big.Int)
		num.DivMod(num, big.NewInt(base), mod)
		encoded = append(encoded, charset[mod.Int64()])
	}
	for len(encoded) < cap(encoded) {
		encoded = append(encoded, '0')
	}
	reverse(encoded)
	return string(encoded)
}
func decodeBase62(s string) (*big.Int, error) {
	result := new(big.Int)
	for _, c := range s {
		val := strings.IndexByte(string(charset), byte(c))
		if val == -1 {
			return nil, ErrInvalidFormat
		}
		result.Mul(result, big.NewInt(base))
		result.Add(result, big.NewInt(int64(val)))
	}
	return result, nil
}

func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}
