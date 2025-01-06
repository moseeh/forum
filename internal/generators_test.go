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

func TestCompareHash(t *testing.T) {
	hash, _ := HashPassword("correctpassword")
	if !CompareHash(hash, "correctpassword") {
		t.Error("Correct password comparison failed")
	}

	if CompareHash(hash, "wrongpassword") {
		t.Error("Incorrect password comparison passed")
	}

	emptyHash, _ := HashPassword("")
	if CompareHash(emptyHash, "") {
		t.Error("Empty password comparison failed")
	}
}
