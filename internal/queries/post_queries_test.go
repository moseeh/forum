package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func dbConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	return db, err
}
func TestGetAllCategories(t *testing.T) {

	db, err := dbConnection()
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	// Create the CATEGORIES table
	createTableQuery := `
		CREATE TABLE CATEGORIES (
			category_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT
		);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %s", err)
	}

	// Insert sample categories
	insertCategoryQuery := `
		INSERT INTO CATEGORIES (category_id, name, description)
		VALUES
			('1', 'Technology', 'Posts about technology'),
			('2', 'Lifestyle', 'Posts about lifestyle'),
			('3', 'Education', 'Posts about education');`
	_, err = db.Exec(insertCategoryQuery)
	if err != nil {
		t.Fatalf("failed to insert categories: %s", err)
	}

	// Initialize the UserModel
	userModel := &UserModel{DB: db}

	// Call GetAllCategories
	categories, err := userModel.GetAllCategories()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Validate the results
	expectedCount := 3
	if len(categories) != expectedCount {
		t.Errorf("unexpected number of categories: got %d, want %d", len(categories), expectedCount)
	}

	expectedCategories := []Category{
		{"1", "Technology", "Posts about technology"},
		{"2", "Lifestyle", "Posts about lifestyle"},
		{"3", "Education", "Posts about education"},
	}

	for i, cat := range categories {
		if cat != expectedCategories[i] {
			t.Errorf("unexpected category at index %d: got %+v, want %+v", i, cat, expectedCategories[i])
		}
	}
}

func TestInsertPost(t *testing.T) {
	db, err := dbConnection()
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	// Create the POSTS and USERS tables
	createTablesQuery := `
		CREATE TABLE USERS (
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);

		CREATE TABLE POSTS (
			post_id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author_id TEXT NOT NULL,
			image_url Text,
			FOREIGN KEY (author_id) REFERENCES USERS (user_id)
		);`
	_, err = db.Exec(createTablesQuery)
	if err != nil {
		t.Fatalf("failed to create tables: %s", err)
	}

	// Insert a sample user
	insertUserQuery := `
		INSERT INTO USERS (user_id, username, email, password)
		VALUES ('1', 'bob', 'bobthebuilder@example.com', 'hashedpassword123');`
	_, err = db.Exec(insertUserQuery)
	if err != nil {
		t.Fatalf("failed to insert test user: %s", err)
	}

	// Initialize the UserModel
	userModel := &UserModel{DB: db}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %s", err)
	}

	// Call InsertPost
	postID := "101"
	title := "Test Post"
	content := "This is a test post content."
	userID := "1"
	imagelink := ""

	err = userModel.InsertPost(tx, postID, title, content, userID, imagelink)
	if err != nil {
		t.Errorf("unexpected error during InsertPost: %s", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		t.Fatalf("failed to commit transaction: %s", err)
	}

	// Validate the inserted post
	var insertedPostID, insertedTitle, insertedContent, insertedAuthorID string
	query := `SELECT post_id, title, content, author_id FROM POSTS WHERE post_id = ?`
	err = db.QueryRow(query, postID).Scan(&insertedPostID, &insertedTitle, &insertedContent, &insertedAuthorID)
	if err != nil {
		t.Errorf("failed to retrieve inserted post: %s", err)
	}

	if insertedPostID != postID || insertedTitle != title || insertedContent != content || insertedAuthorID != userID {
		t.Errorf("unexpected post values: got (%s, %s, %s, %s), want (%s, %s, %s, %s)",
			insertedPostID, insertedTitle, insertedContent, insertedAuthorID,
			postID, title, content, userID)
	}
}
