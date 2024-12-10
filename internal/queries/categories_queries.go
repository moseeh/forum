package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
)

var categories = []Category{
	{Name: "Technology", Description: "All about the latest tech trends and innovations"},
	{Name: "Sports", Description: "Updates and news about various sports"},
	{Name: "Music", Description: "Genres, artists, and music events"},
	{Name: "Food", Description: "Delicious recipes and culinary delights"},
	{Name: "Science", Description: "Discoveries and research in science"},
	{Name: "AI", Description: "Artificial Intelligence advancements and discussions"},
	{Name: "Travel", Description: "Travel guides and destinations"},
	{Name: "Health", Description: "Health tips and medical advice"},
	{Name: "Finance", Description: "Financial advice and economic news"},
}

// Function to insert categories into the database
func InsertCategories(db *sql.DB) {
	// Use INSERT OR IGNORE to avoid duplicate entries based on the UNIQUE constraint
	stmt, err := db.Prepare(`INSERT OR IGNORE INTO CATEGORIES (category_id, name, description) VALUES (?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, category := range categories {
		id, err := uuid.NewV4()
		if err != nil {
			log.Fatal("Error generating UUID:", err)
		}

		_, err = stmt.Exec(id.String(), category.Name, category.Description)
		if err != nil {
			log.Printf("Error inserting category %s: %v\n", category.Name, err)
		} else {
			fmt.Printf("Inserted or skipped category: %s\n", category.Name)
		}
	}
}
