package identifier

import (
	"bytes"
	"math/big"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	id := ID("prefix").New()

	if len(id) != timestampLength+1+randomLength {
		t.Errorf("New() returned ID of incorrect length. Got %d, want %d", len(id), timestampLength+1+randomLength)
	}

	parts := strings.Split(string(id), string(separatorChar))
	if len(parts) != 2 {
		t.Errorf("New() returned ID with incorrect format. Got %s", id)
	}

	if len(parts[0]) != timestampLength || len(parts[1]) != randomLength {
		t.Errorf("New() returned ID with incorrect part lengths. Got %d and %d, want %d and %d",
			len(parts[0]), len(parts[1]), timestampLength, randomLength)
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid ID", "1234567890-123456789", false},
		{"Invalid length", "123456789-12345678", true},
		{"Invalid separator", "1234567890123456789", true},
		{"Invalid characters", "123456789O-123456789", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ID("").FromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFromBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  ID
	}{
		{"Valid bytes", []byte("1234567890-123456789"), "1234567890-123456789"},
		{"Empty bytes", []byte{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ID("").FromBytes(tt.input)
			if err != nil {
				t.Errorf("FromBytes() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("FromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		id      ID
		want    bool
		wantErr bool
	}{
		{"Valid ID", "1234567890-123456789", true, false},
		{"Invalid length", "123456789-12345678", false, true},
		{"Invalid separator", "1234567890123456789", false, true},
		{"Invalid characters", "123456789O-123456789", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.id.Validate(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		name  string
		input *big.Int
		want  string
	}{
		{"Zero", big.NewInt(0), "0000000000"},
		{"Positive number", big.NewInt(12345), "0000012345"},
		{"Large number", new(big.Int).SetBytes([]byte{255, 255, 255, 255}), "2LKcb1ZMR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeBase62(tt.input); got != tt.want {
				t.Errorf("encodeBase62() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeBase62(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *big.Int
		wantErr bool
	}{
		{"Zero", "0000000000", big.NewInt(0), false},
		{"Positive number", "0000012345", big.NewInt(12345), false},
		{"Large number", "2LKcb1ZMR", new(big.Int).SetBytes([]byte{255, 255, 255, 255}), false},
		{"Invalid character", "000001234O", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeBase62(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeBase62() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got.Cmp(tt.want) != 0 {
				t.Errorf("decodeBase62() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{"Empty slice", []byte{}, []byte{}},
		{"Single element", []byte{1}, []byte{1}},
		{"Multiple elements", []byte{1, 2, 3, 4, 5}, []byte{5, 4, 3, 2, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]byte, len(tt.input))
			copy(input, tt.input)
			reverse(input)
			if !bytes.Equal(input, tt.want) {
				t.Errorf("reverse() = %v, want %v", input, tt.want)
			}
		})
	}
}
