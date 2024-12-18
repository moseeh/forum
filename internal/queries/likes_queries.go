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

