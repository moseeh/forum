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
	if !CompareHash(emptyHash, "") {
		t.Error("Empty password comparison failed")
	}
}

func TestUUIDGen(t *testing.T) {
	uuid1 := UUIDGen()
	uuid2 := UUIDGen()

	if uuid1 == "" || uuid2 == "" {
		t.Error("UUID generation produced an empty string")
	}

	if uuid1 == uuid2 {
		t.Error("UUID generation produced the same UUID twice, which should be rare")
	}
}

func TestTokenGen(t *testing.T) {
	// Test token generation with different lengths
	token1 := TokenGen(16)
	token2 := TokenGen(32)

	if len(token1) != 24 || len(token2) != 44 { // Note: base64 encoding increases the size
		t.Errorf("Token length incorrect. Expected 24 for 16 bytes input, got %d", len(token1))
		t.Errorf("Token length incorrect. Expected 44 for 32 bytes input, got %d", len(token2))
	}

	// Check if tokens are different
	if token1 == token2 {
		t.Error("Two tokens with different length requirements are identical")
	}
}
