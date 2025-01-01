package db

import (
	"database/sql"
	"log"

	"github.com/gofrs/uuid"
)

type CategoryCount struct {
	CategoryID   string
	Count        int
	CategoryName string
}

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
		}
	}
}

func (m *UserModel) TrendingCount() ([]CategoryCount, error) {
	rows, err := m.DB.Query(` SELECT pc.category_id, COUNT(*) AS count, c.name FROM POST_CATEGORIES pc JOIN CATEGORIES c ON pc.category_id = c.category_id GROUP BY pc.category_id `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryCounts []CategoryCount

	for rows.Next() {
		var categoryCount CategoryCount
		if err := rows.Scan(&categoryCount.CategoryID, &categoryCount.Count, &categoryCount.CategoryName); err != nil {
			return nil, err
		}
		categoryCounts = append(categoryCounts, categoryCount)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	//sort
	for i := 0; i < len(categoryCounts)-1; i++ {
		if categoryCounts[i+1].Count > categoryCounts[i].Count {
			categoryCounts[i], categoryCounts[i+1] = categoryCounts[i+1], categoryCounts[i]
		}
	}

	return categoryCounts, nil
}
