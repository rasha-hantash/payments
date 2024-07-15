package identifier

import (
	"time"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		prefix ID
	}{
		{"With acct_ prefix", "acct_"},
		{"With user_ prefix", "user_"},
		{"With txn_ prefix", "txn_"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.prefix.New()

			// Check that the ID starts with the correct prefix
			if !hasPrefix(string(id), string(tt.prefix)) {
				t.Errorf("New() = %v, does not start with prefix %v", id, tt.prefix)
			}

			// Check that the total length of the ID is correct
			if len(id) != idLength {
				t.Errorf("New() = %v, length = %d, want length = %d", id, len(id), idLength)
			}

			// Extract the timestamp and random part from the ID
			timestamp := string(id[len(tt.prefix) : len(tt.prefix)+12])
			randomPart := string(id[len(tt.prefix)+12:])

			// Validate the timestamp part
			_, err := time.Parse("060102150405", timestamp)
			if err != nil {
				t.Errorf("New() = %v, invalid timestamp = %v", id, timestamp)
			}

			// Validate the random part
			for idx, char := range randomPart {
				if idx%2 == 0 {
					if !contains(alphabet, char) {
						t.Errorf("New() = %v, invalid random part character at %d = %c", id, idx, char)
					}
				} else {
					if !contains(digits, char) {
						t.Errorf("New() = %v, invalid random part character at %d = %c", id, idx, char)
					}
				}
			}
		})
	}
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}



func TestFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid ID", "acct_240714212559C7E", false},
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
		wantErr bool
	}{
		{"Valid bytes", []byte("acct_240714212559C7E"), "acct_240714212559C7E", false},
		{"Empty bytes", []byte{}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ID("acct_").FromBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromBytes() error = %v, wantErr %v", err, tt.wantErr)
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
		name string
		id   ID
		want bool
	}{
		{"Valid ID", "acct_240714212559C7E", true},
		{"Invalid length", "acct_240714212559A1", false},
		{"Invalid characters", "acct_240714212559A1B@", false},
		{"Invalid timestamp", "acct_991231123059A1B2", false},
		{"Missing underscore", "acct240714212559A1B2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.id.Validate()
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
