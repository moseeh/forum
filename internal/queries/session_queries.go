package db

import (
	"database/sql"
	"fmt"
	"log"
)

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

func (m *UserModel) ValidateSession(sessionToken string) (bool, error) {
	query := `
        SELECT user_id 
        FROM TOKENS 
        WHERE session_token = ? 
        AND expires_at > datetime('now')
        AND user_id IS NOT NULL`

	var userID string
	err := m.DB.QueryRow(query, sessionToken).Scan(&userID)
	if err != nil {

		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("database error while validating session: %v", err)
	}

	updateQuery := `
        UPDATE TOKENS 
        SET expires_at = datetime('now', '+2 hours')
        WHERE session_token = ?`

	_, err = m.DB.Exec(updateQuery, sessionToken)
	if err != nil {
		log.Printf("Failed to update session expiry: %v", err)
	}
	return true, nil
}

func (m *UserModel) DeleteUserSessions(userID string) error {
	const DELETE_USER_SESSIONS = `
	DELETE FROM TOKENS 
	WHERE user_id = ?`

	_, err := m.DB.Exec(DELETE_USER_SESSIONS, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %v", err)
	}
	return nil
}

func (m *UserModel) DeleteSession(sessionToken string) error {
    const DELETE_SESSION = `
        DELETE FROM TOKENS 
        WHERE session_token = ?`
    
    _, err := m.DB.Exec(DELETE_SESSION, sessionToken)
    if err != nil {
        return fmt.Errorf("failed to delete session: %v", err)
    }
    return nil
}
