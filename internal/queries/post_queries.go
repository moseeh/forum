package db

import (
	"fmt"
	"time"
)

type Post struct {
	PostID    string    `json:"post_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"upadted_at"`
	AuthorID  string    `json:"author_id"`
	Username  string    `json:"username"`
}

func (m *UserModel) InsertPost(post_id, title, content, user_id string) error {
	const POST_INSERT string = "INSERT INTO POSTS (post_id, title, content, author_id) VALUES (?,?,?,?)"
	stmt, err := m.DB.Prepare(POST_INSERT)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(post_id, title, content, user_id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *UserModel) GetALlPosts() ([]Post, error) {
	const GET_ALL_POSTS = `
        SELECT p.post_id, p.title, p.content, p.created_at, p.updated_at, 
               p.author_id, u.username
        FROM POSTS p
        LEFT JOIN USERS u ON p.author_id = u.user_id
        ORDER BY p.created_at DESC`
	stmt, err := m.DB.Prepare(GET_ALL_POSTS)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&post.PostID,
			&post.Title,
			&post.Content,
			&createdAtStr,
			&updatedAtStr,
			&post.AuthorID,
			&post.Username,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return posts, nil
}
