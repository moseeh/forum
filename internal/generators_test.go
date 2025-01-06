package internal

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	// Test case for hashing a password
	password := "testpassword"
	hash, err := HashPassword(password)

	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("Generated hash is empty")
	}

	if !CompareHash(hash, password) {
		t.Error("Hash doesn't match original password")
	}

	emptyHash, err := HashPassword("")
	if err != nil {
		t.Errorf("HashPassword with empty string failed: %v", err)
	}
	if emptyHash == "" {
		t.Error("Empty password hash is empty")
	}
}
