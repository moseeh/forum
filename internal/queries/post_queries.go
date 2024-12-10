package db

import (
	"database/sql"
	"time"
)

type Post struct {
	PostID        string    `json:"post_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"upadted_at"`
	AuthorID      string    `json:"author_id"`
	Username      string    `json:"username"`
	Categories    []Category
	LikesCount    int
	CommentsCount int
	IsLiked       bool
}

type Category struct {
	CategoryID  string
	Name        string
	Description string
}

// GetAllCategories retrieves all available categories
func (u *UserModel) GetAllCategories() ([]Category, error) {
	query := `SELECT category_id, name, description FROM CATEGORIES`
	rows, err := u.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		err := rows.Scan(&cat.CategoryID, &cat.Name, &cat.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// InsertPost creates a new post within a transaction
func (u *UserModel) InsertPost(tx *sql.Tx, postID, title, content, userID string) error {
	query := `INSERT INTO POSTS (post_id, title, content, author_id) VALUES (?, ?, ?, ?)`
	_, err := tx.Exec(query, postID, title, content, userID)
	return err
}

// InsertPostCategory links a post to a category
func (u *UserModel) InsertPostCategory(tx *sql.Tx, postID, categoryID string) error {
	query := `INSERT INTO POST_CATEGORIES (post_id, category_id) VALUES (?, ?)`
	_, err := tx.Exec(query, postID, categoryID)
	return err
}

// GetAllPosts retrieves all posts with their categories, likes, and comments
func (u *UserModel) GetAllPosts(currentUserID string) ([]Post, error) {
	query := `
        SELECT 
            p.post_id, p.title, p.content, p.created_at, p.updated_at,
            u.username, p.author_id, p.likes_count, p.comments_count,
            CASE 
                WHEN l.user_id IS NOT NULL THEN 1 
                ELSE 0 
            END as is_liked
        FROM POSTS p
        JOIN USERS u ON p.author_id = u.user_id
        LEFT JOIN LIKES l ON p.post_id = l.post_id AND l.user_id = ?
        ORDER BY p.created_at DESC`

	rows, err := u.DB.Query(query, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.PostID, &post.Title, &post.Content,
			&post.CreatedAt, &post.UpdatedAt, &post.Username,
			&post.AuthorID, &post.LikesCount, &post.CommentsCount,
			&post.IsLiked,
		)
		if err != nil {
			return nil, err
		}

		// Get categories for this post
		categories, err := u.GetPostCategories(post.PostID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}
	return posts, nil
}

// GetPostCategories retrieves all categories for a specific post
func (u *UserModel) GetPostCategories(postID string) ([]Category, error) {
	query := `
        SELECT c.category_id, c.name, c.description
        FROM CATEGORIES c
        JOIN POST_CATEGORIES pc ON c.category_id = pc.category_id
        WHERE pc.post_id = ?`

	rows, err := u.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		err := rows.Scan(&cat.CategoryID, &cat.Name, &cat.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}
