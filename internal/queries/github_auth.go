package db

import "database/sql"

type User struct {
	UserID       string
	Username     string
	AuthProvider string
}

func (m *UserModel) GetUserbYUsername(username string) (*User, error) {
	var user User

	query := `SELECT user_id, username, auth_provider FROM USERS WHERE username = ?`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.UserID, &user.Username, &user.AuthProvider)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
