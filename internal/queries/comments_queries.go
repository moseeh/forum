package db

import "time"

type Comment struct {
	Comment_ID    string
	Content       string
	Username      string
	CreatedAt     time.Time
	LikesCount    int
	DislikesCount int
	IsLiked       bool
	IsDisliked    bool
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

func (m *UserModel) GetPostComments(postID, currentUserID string) ([]Comment, error) {
	query := `
    SELECT
        c.comment_id,
        c.content,
        u.username,
        c.created_at,
        c.likes_count,
        c.dislikes_count,
        CASE
            WHEN cl.user_id IS NOT NULL THEN 1
            ELSE 0
        END as is_liked,
        CASE
            WHEN cd.user_id IS NOT NULL THEN 1
            ELSE 0
        END as is_disliked
    FROM COMMENTS c
    JOIN USERS u ON c.user_id = u.user_id
    LEFT JOIN COMMENT_LIKES cl ON c.comment_id = cl.comment_id AND cl.user_id = ?
    LEFT JOIN COMMENT_DISLIKES cd ON c.comment_id = cd.comment_id AND cd.user_id = ?
    WHERE c.post_id = ?
    ORDER BY c.created_at DESC`

	rows, err := m.DB.Query(query, currentUserID, currentUserID, postID)
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
			&comment.LikesCount,
			&comment.DislikesCount,
			&comment.IsLiked,
			&comment.IsDisliked,
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
