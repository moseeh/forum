package db

import "time"

type Comment struct {
	Comment_ID string
	Content    string
	Username   string
	CreatedAt  time.Time
}

func (m *UserModel) InsertComment(comment_id, post_id, user_id, content string) error {
	query := `INSERT INTO COMMENTS (comment_id, post_id, user_id, content) VALUES (?,?,?,?)`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(comment_id, post_id, user_id, content)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) GetPostComments(postID string) ([]Comment, error) {
	query := `
    SELECT 
        c.comment_id,
        c.content,
        u.username,
        c.created_at
    FROM COMMENTS c
    JOIN USERS u ON c.user_id = u.user_id
    WHERE c.post_id = ?
    ORDER BY c.created_at DESC`

	rows, err := m.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.Comment_ID,
			&comment.Content,
			&comment.Username,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
