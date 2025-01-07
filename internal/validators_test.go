package internal

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	// Valid usernames
	validUsernames := []string{"a", "A", "user_name", "UserName123", "a_b_c"}
	for _, username := range validUsernames {
		if err := ValidateUsername(username); err != nil {
			t.Errorf("Expected valid username %s, but got error: %v", username, err)
		}
	}

	// Invalid usernames
	invalidUsernames := []string{
		"",           // empty string
		"1user",      // starts with number
		"user name",  // contains space
		"user@name",  // contains @
		"user-name",  // contains -
		"user..name", // contains ..
		"__user",     // starts with double underscore
		"user__name", // contains double underscore
		"admin",      // reserved
		"root",       // reserved
		"system",     // reserved
		"test",       // reserved
		"null",       // reserved
		"localhost",  // reserved
		"void",       // reserved
		"guest",      // reserved
	}

	for _, username := range invalidUsernames {
		err := ValidateUsername(username)
		if err == nil {
			t.Errorf("Expected invalid username %s to return an error, but got nil", username)
		} else {
			t.Logf("Successfully detected invalid username %s with error: %v", username, err)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	validEmails := []string{
		"user@example.com",
		"user.name+tag@example.co.uk",
		"user-name@example.org",
	}
	for _, email := range validEmails {
		if !ValidateEmail(email) {
			t.Errorf("Expected %s to be a valid email", email)
		}
	}

	invalidEmails := []string{
		"user@example",
		"user@.com",
		"user@example.",
		"@example.com",
		"user name@example.com",
		"user@example.c",
	}
	for _, email := range invalidEmails {
		if ValidateEmail(email) {
			t.Errorf("Expected %s to be an invalid email", email)
		}
	}
}
