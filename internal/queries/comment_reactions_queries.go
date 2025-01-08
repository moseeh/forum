package db

func (m *UserModel) UserLikeOnCommentExists(commentID, userID string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM COMMENT_LIKES WHERE comment_id = ? AND user_id = ?)`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var exists bool

	err = stmt.QueryRow(commentID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (m *UserModel) InsertCommentLike(likeID, commentID, userID string) error {
	query := `INSERT INTO COMMENT_LIKES (like_id, comment_id, user_id) VALUES (?,?,?)`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(likeID, commentID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) DeleteCommentLike(commentID, userID string) error {
	query := `DELETE FROM COMMENT_LIKES WHERE comment_id = ? AND user_id = ?`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(commentID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UserDislikeOnCommentExists(commentID, userID string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM COMMENT_DISLIKES WHERE comment_id = ? AND user_id = ?)`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var exists bool

	err = stmt.QueryRow(commentID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (m *UserModel) InsertCommentDislike(likeID, commentID, userID string) error {
	query := `INSERT INTO COMMENT_DISLIKES (dislike_id, comment_id, user_id) VALUES (?,?,?)`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(likeID, commentID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) DeleteCommentDislike(commentID, userID string) error {
	query := `DELETE FROM COMMENT_DISLIKES WHERE comment_id = ? AND user_id = ?`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(commentID, userID)
	if err != nil {
		return err
	}
	return nil
}
