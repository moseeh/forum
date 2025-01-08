package db

func (u *UserModel) GetLikedPosts(currentUserID string) ([]Post, error) {
	query := `
	SELECT 
		p.post_id, p.title, p.content, p.created_at, p.updated_at,
		u.username, p.author_id, p.likes_count, p.dislikes_count, p.comments_count,
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
	WHERE l.user_id IS NOT NULL OR d.user_id IS NOT NULL  
	ORDER BY p.created_at DESC`

	rows, err := u.DB.Query(query, currentUserID, currentUserID)
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
			&post.AuthorID, &post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.IsLiked, &post.IsDisliked,
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

func (u *UserModel) GetCreatedPosts(currentUserID string) ([]Post, error) {
	query := `
	SELECT 
		p.post_id, p.title, p.content, p.created_at, p.updated_at,
		u.username, p.author_id, p.likes_count, p.dislikes_count, p.comments_count,
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
	WHERE p.author_id = ?
	ORDER BY p.created_at DESC`

	rows, err := u.DB.Query(query, currentUserID, currentUserID, currentUserID)
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
			&post.AuthorID, &post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.IsLiked, &post.IsDisliked,
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
