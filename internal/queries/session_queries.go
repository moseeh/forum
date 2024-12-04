package db

func (m *UserModel) NewSession(user_id, session_token, csrf_token, expires_at string) error {
	const INSERT_TOKENS string = "INSERT INTO TOKENS (user_id, session_token,csrf_token,expires_at) VALUES (?, ?, ?, ?)"
	stmt, err := m.DB.Prepare(INSERT_TOKENS)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user_id, session_token, csrf_token, expires_at)
	if err != nil {
		return err
	}
	return nil
}
