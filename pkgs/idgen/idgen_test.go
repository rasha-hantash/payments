package idgen

import (
	"regexp"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	handler := CustomIDHandler{}

	for i := 0; i < 1000; i++ {
		id, err := handler.New()
		if err != nil {
			t.Errorf("New() returned an error: %v", err)
		}

		if len(id) != idLength {
			t.Errorf("New() returned ID of length %d, want %d", len(id), idLength)
		}

		valid, err := handler.Validate(id)
		if err != nil || !valid {
			t.Errorf("New() returned invalid ID: %s", id)
		}
	}
}

func TestFromString(t *testing.T) {
	handler := CustomIDHandler{}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid ID", "01234567890123456789", false},
		{"Too short", "0123456789", true},
		{"Too long", "012345678901234567890", true},
		{"Invalid characters", "01234567890123456789!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.FromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFromBytes(t *testing.T) {
	handler := CustomIDHandler{}

	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{"Valid ID", []byte("01234567890123456789"), false},
		{"Too short", []byte("0123456789"), true},
		{"Too long", []byte("012345678901234567890"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.FromBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	handler := CustomIDHandler{}

	tests := []struct {
		name    string
		input   CustomID
		want    bool
		wantErr bool
	}{
		{"Valid ID", "01234567890123456789", true, false},
		{"Too short", "0123456789", false, true},
		{"Too long", "012345678901234567890", false, true},
		{"Invalid characters", "01234567890123456789!", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.Validate(tt.input)
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

func TestMonotonicity(t *testing.T) {
	handler := CustomIDHandler{}
	var prevID CustomID

	for i := 0; i < 1000; i++ {
		id, err := handler.New()
		if err != nil {
			t.Fatalf("New() returned an error: %v", err)
		}

		if i > 0 && id <= prevID {
			t.Errorf("New() returned non-monotonic ID: prev = %s, current = %s", prevID, id)
		}

		prevID = id
		time.Sleep(time.Millisecond) // Ensure different timestamps
	}
}

func TestEncoding(t *testing.T) {
	// Test that the encoding uses the correct character set
	validChars := "^[0-9A-HJ-NP-TV-Z]+$"
	re := regexp.MustCompile(validChars)

	handler := CustomIDHandler{}

	for i := 0; i < 1000; i++ {
		id, err := handler.New()
		if err != nil {
			t.Fatalf("New() returned an error: %v", err)
		}

		if !re.MatchString(string(id)) {
			t.Errorf("New() returned ID with invalid characters: %s", id)
		}
	}
}

func TestTimestampExtraction(t *testing.T) {
	handler := CustomIDHandler{}

	id, err := handler.New()
	if err != nil {
		t.Fatalf("New() returned an error: %v", err)
	}

	// Decode the first 8 characters (6 bytes) to get the timestamp
	decoded, err := encoding.DecodeString(string(id[:8]))
	if err != nil {
		t.Fatalf("Failed to decode ID: %v", err)
	}

	// Extract timestamp
	var timestamp int64
	for i := 0; i < 6; i++ {
		timestamp = timestamp<<8 | int64(decoded[i])
	}

	// Check if the extracted timestamp is close to the current time
	now := time.Now().UnixNano() / 1e6
	if now-timestamp > 1000 || timestamp-now > 1000 {
		t.Errorf("Extracted timestamp %d is not close to current time %d", timestamp, now)
	}
}

func BenchmarkNew(b *testing.B) {
	handler := CustomIDHandler{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler.New()
	}
}

func BenchmarkFromString(b *testing.B) {
	handler := CustomIDHandler{}
	id, _ := handler.New()
	idStr := string(id)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler.FromString(idStr)
	}
}

func BenchmarkFromBytes(b *testing.B) {
	handler := CustomIDHandler{}
	id, _ := handler.New()
	idBytes := []byte(id)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler.FromBytes(idBytes)
	}
}

func BenchmarkValidate(b *testing.B) {
	handler := CustomIDHandler{}
	id, _ := handler.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler.Validate(id)
	}
}
