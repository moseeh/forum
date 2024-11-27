package db

import (
	"fmt"
	"log"
)

// Query Strings
const USER_INSERT string = "INSERT INTO USERS (user_id, username, email, password) VALUES (?,?,?,?);"
const USER_EXISTS string = "SELECT COUNT(1) FROM USERS WHERE username = ? OR email = ?;"
const PASSWORD_HASH string = "SELECT password FROM USERS WHERE email = ?;"
const CHECK_USER string = "SELECT COUNT(1) FROM USERS WHERE email = ?;"
const INSERT_TOKENS string = "INSERT INTO TOKENS (user_id, session_token,csrf_token,expires_at) VALUES (?, ?, ?, ?)"
const UPDATE_TOKENS string = "UPDATE TOKENS SET session_token = ?,csrf_token = ?,  expires_at = ? WHERE user_id = ?"
const USER_ID string = "SELECT user_id FROM USERS WHERE email = ?"
const EMAIL_EXIST string = "SELECT COUNT(1) FROM USERS WHERE email = ?;"
const GET_SESSION_TOKEN string = "SELECT session_token FROM TOKENS WHERE user_id = ?"

func (conn ForumDB) Insert(query string, values ...interface{}) error {
	stmt, err := conn.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	res, err := stmt.Exec(values...)
	if err != nil {
		log.Println(res)
		return err
	}
	return nil
}

func (conn ForumDB) Exists(query string, values ...interface{}) (bool, error) {
	stmt, err := conn.db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(values...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (conn ForumDB) GetStringValue(query string, values ...interface{}) (string, error) {
	stmt, err := conn.db.Prepare(query)
	if err != nil {
		return "", nil
	}
	defer stmt.Close()

	var passwordHash string

	if err = stmt.QueryRow(values...).Scan(&passwordHash); err != nil {
		return "", err
	}
	return passwordHash, nil
}
