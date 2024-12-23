package db

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
