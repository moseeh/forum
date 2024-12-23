package db

func (m *UserModel) UserDislikeOnPostExists(postID, userID string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM DISLIKES WHERE post_id = ? AND user_id = ?)`
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

func (m *UserModel) InsertDislike(likeID, postID, userID string) error {
	query := `INSERT INTO DISLIKES (dislike_id, post_id, user_id) VALUES (?,?,?)`
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

func (m *UserModel) DeleteDislike(postID, userID string) error {
	query := `DELETE FROM DISLIKES WHERE post_id = ? AND user_id = ?`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(postID, userID)
	if err != nil {
		return err
	}
	return nil
}
