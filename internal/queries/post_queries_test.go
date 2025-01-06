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
