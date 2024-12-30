package db

import (
	"database/sql"
	"fmt"
	"log"
)

func (m *UserModel) NewSession(user_id, session_token, csrf_token, expires_at string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// First, delete any existing sessions for this user
	const DELETE_EXISTING = `
        DELETE FROM TOKENS 
        WHERE user_id = ?`

	_, err = tx.Exec(DELETE_EXISTING, user_id)
	if err != nil {
		return fmt.Errorf("failed to delete existing sessions: %v", err)
	}

	// Then create the new session
	const INSERT_TOKENS = `
        INSERT INTO TOKENS (user_id, session_token, csrf_token, expires_at) 
        VALUES (?, ?, ?, ?)`

	_, err = tx.Exec(INSERT_TOKENS, user_id, session_token, csrf_token, expires_at)
	if err != nil {
		return fmt.Errorf("failed to insert new session: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
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
