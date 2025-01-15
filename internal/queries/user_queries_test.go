package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInsertUser(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	createTableQuery := `
		CREATE TABLE USERS (
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			avatar_url TEXT,
            auth_provider TEXT NOT NULL
		);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %s", err)
	}

	userModel := &UserModel{DB: db}

	// Define test inputs
	id := "abcd"
	username := "aaochieng"
	email := "aaochieng@example.com"
	password := "securepassword"
	authProvider := "traditional"
	avatarUrl := ""

	err = userModel.InsertUser(id, username, email, password, authProvider, avatarUrl)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	var (
		gotID       string
		gotUsername string
		gotEmail    string
		gotPassword string
	)

	query := "SELECT user_id, username, email, password FROM USERS WHERE user_id = ?"
	row := db.QueryRow(query, id)
	err = row.Scan(&gotID, &gotUsername, &gotEmail, &gotPassword)
	if err != nil {
		t.Fatalf("failed to query inserted user: %s", err)
	}

	if gotID != id || gotUsername != username || gotEmail != email || gotPassword != password {
		t.Errorf("inserted user data does not match: got (%s, %s, %s, %s), want (%s, %s, %s, %s)",
			gotID, gotUsername, gotEmail, gotPassword, id, username, email, password)
	}
}

func TestUserExists(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	createTableQuery := `
		CREATE TABLE USERS (
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %s", err)
	}

	userModel := &UserModel{DB: db}

	insertUserQuery := `
		INSERT INTO USERS (user_id, username, email, password)
		VALUES ('1', 'squidward', 'squidward@example.com', 'securepassword');`
	_, err = db.Exec(insertUserQuery)
	if err != nil {
		t.Fatalf("failed to insert test user: %s", err)
	}

	tests := []struct {
		email       string
		expected    bool
		shouldError bool
	}{
		{"squidward@example.com", true, false},    // User exists
		{"nonexistent@example.com", false, false}, // User does not exist
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			// Call the method under test
			exists, err := userModel.UserEmailExists(test.email)

			// Check for unexpected errors
			if test.shouldError && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !test.shouldError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			// Verify the result
			if exists != test.expected {
				t.Errorf("unexpected result: got %v, want %v", exists, test.expected)
			}
		})
	}
}

func TestGetPassword(t *testing.T) {
	// Set up an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	// Create the USERS table
	createTableQuery := `
		CREATE TABLE USERS (
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %s", err)
	}

	// Initialize the UserModel
	userModel := &UserModel{DB: db}

	// Insert a test user
	insertUserQuery := `
		INSERT INTO USERS (user_id, username, email, password)
		VALUES ('1', 'testuser', 'mrcrabs@example.com', 'hashedpassword123');`
	_, err = db.Exec(insertUserQuery)
	if err != nil {
		t.Fatalf("failed to insert test user: %s", err)
	}

	// Test cases
	tests := []struct {
		email       string
		expected    string
		shouldError bool
	}{
		{"mrcrabs@example.com", "hashedpassword123", false}, // Valid email
		{"mrspuff@example.com", "", true},                   // Nonexistent email
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			// Call the method under test
			password, err := userModel.GetPassword(test.email)

			// Check for unexpected errors
			if test.shouldError && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !test.shouldError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			// Verify the result
			if password != test.expected {
				t.Errorf("unexpected result: got %v, want %v", password, test.expected)
			}
		})
	}
}

func TestGetUsername(t *testing.T) {
	// Set up an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	// Create the USERS table
	createTableQuery := `
		CREATE TABLE USERS (
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %s", err)
	}

	// Initialize the UserModel
	userModel := &UserModel{DB: db}

	// Insert a test user
	insertUserQuery := `
		INSERT INTO USERS (user_id, username, email, password)
		VALUES ('1', 'ben10', 'ben10@example.com', 'hashedpassword123');`
	_, err = db.Exec(insertUserQuery)
	if err != nil {
		t.Fatalf("failed to insert test user: %s", err)
	}

	// Test GetUsername
	t.Run("Valid email", func(t *testing.T) {
		userID, username, err := userModel.GetUsername("ben10@example.com")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if userID != "1" || username != "ben10" {
			t.Errorf("unexpected result: got (%s, %s), want (1, ben10)", userID, username)
		}
	})

	t.Run("Nonexistent email", func(t *testing.T) {
		_, _, err := userModel.GetUsername("nonexistent@example.com")
		if err == nil {
			t.Error("expected an error but got none")
		}
	})
}

func TestGetUserID(t *testing.T) {
	// Set up an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	// Create the USERS table
	createTableQuery := `
		CREATE TABLE USERS (
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %s", err)
	}

	// Initialize the UserModel
	userModel := &UserModel{DB: db}

	// Insert a test user
	insertUserQuery := `
		INSERT INTO USERS (user_id, username, email, password)
		VALUES ('1', 'randomuser', 'randuser@example.com', 'hashedpassword123');`
	_, err = db.Exec(insertUserQuery)
	if err != nil {
		t.Fatalf("failed to insert test user: %s", err)
	}

	// Test GetUserID
	t.Run("Valid username", func(t *testing.T) {
		userID, err := userModel.GetUserID("randomuser")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if userID != "1" {
			t.Errorf("unexpected result: got %s, want 1", userID)
		}
	})

	t.Run("Nonexistent username", func(t *testing.T) {
		_, err := userModel.GetUserID("nonexistent")
		if err == nil {
			t.Error("expected an error but got none")
		}
	})
}

func TestInsertPostCategory(t *testing.T) {
	db, err := dbConnection()
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	// Create the POSTS, CATEGORIES, and POST_CATEGORIES tables
	createTablesQuery := `
		CREATE TABLE POSTS (
			post_id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author_id TEXT NOT NULL
		);

		CREATE TABLE CATEGORIES (
			category_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT
		);

		CREATE TABLE POST_CATEGORIES (
			post_id TEXT NOT NULL,
			category_id TEXT NOT NULL,
			FOREIGN KEY (post_id) REFERENCES POSTS (post_id),
			FOREIGN KEY (category_id) REFERENCES CATEGORIES (category_id)
		);`
	_, err = db.Exec(createTablesQuery)
	if err != nil {
		t.Fatalf("failed to create tables: %s", err)
	}

	// Insert sample data for POSTS and CATEGORIES
	insertSampleDataQuery := `
		INSERT INTO POSTS (post_id, title, content, author_id)
		VALUES ('1', 'Sample Post', 'Sample content', 'author1');

		INSERT INTO CATEGORIES (category_id, name, description)
		VALUES ('cat1', 'Technology', 'Posts about technology');`
	_, err = db.Exec(insertSampleDataQuery)
	if err != nil {
		t.Fatalf("failed to insert sample data: %s", err)
	}

	// Initialize the UserModel
	userModel := &UserModel{DB: db}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %s", err)
	}

	// Call InsertPostCategory
	postID := "1"
	categoryID := "cat1"

	err = userModel.InsertPostCategory(tx, postID, categoryID)
	if err != nil {
		t.Errorf("unexpected error during InsertPostCategory: %s", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		t.Fatalf("failed to commit transaction: %s", err)
	}

	// Validate the inserted post-category mapping
	var insertedPostID, insertedCategoryID string
	query := `SELECT post_id, category_id FROM POST_CATEGORIES WHERE post_id = ? AND category_id = ?`
	err = db.QueryRow(query, postID, categoryID).Scan(&insertedPostID, &insertedCategoryID)
	if err != nil {
		t.Errorf("failed to retrieve inserted post-category mapping: %s", err)
	}

	if insertedPostID != postID || insertedCategoryID != categoryID {
		t.Errorf("unexpected mapping values: got (%s, %s), want (%s, %s)",
			insertedPostID, insertedCategoryID, postID, categoryID)
	}
}
