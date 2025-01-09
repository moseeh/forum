package db

import "database/sql"

func (m *UserModel) GetPostDetails(postID, currentUserID string) (*Post, error) {
	query := `
    SELECT 
        p.post_id, p.title, p.content, p.created_at, p.updated_at,
        u.username, p.author_id, p.image_url, p.likes_count, p.dislikes_count, p.comments_count,
        CASE 
            WHEN l.user_id IS NOT NULL THEN 1
            ELSE 0
        END as is_liked,
        CASE 
            WHEN d.user_id IS NOT NULL THEN 1
            ELSE 0
        END as is_disliked
    FROM POSTS p
    JOIN USERS u ON p.author_id = u.user_id
    LEFT JOIN LIKES l ON p.post_id = l.post_id AND l.user_id = ?
    LEFT JOIN DISLIKES d ON p.post_id = d.post_id AND d.user_id = ?
    WHERE p.post_id = ?`

	var post Post
	err := m.DB.QueryRow(query, currentUserID, currentUserID, postID).Scan(
		&post.PostID, &post.Title, &post.Content,
		&post.CreatedAt, &post.UpdatedAt, &post.Username,
		&post.AuthorID, &post.ImageName, &post.LikesCount, &post.DislikesCount, &post.CommentsCount,
		&post.IsLiked, &post.IsDisliked,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// Get categories for this post
	categories, err := m.GetPostCategories(post.PostID)
	if err != nil {
		return nil, err
	}
	post.Categories = categories

	return &post, nil
}
