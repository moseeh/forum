package db

func (m *UserModel) UserLikeOnPostExists(postID, userID string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM LIKES WHERE post_id = ? AND user_id = ?)`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var exists bool

	err = stmt.QueryRow(postID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (m *UserModel) InsertLike(postID, userID, likeID string) error {
	query := `INSERT INTO LIKES (like_id, post_id, user_id) VALUES (?,?,?)`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(likeID, postID, userID)
	if err != nil {
		return err
	}

	return nil
}
