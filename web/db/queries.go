package db

import (
	"fmt"
	"log"
)

// Query Strings
const USER_INSERT string = "INSERT INTO USERS (user_id, username, email, password) VALUES (?,?,?,?);"
const USER_EXISTS string = "SELECT COUNT(1) FROM USERS WHERE username = ? OR email = ?;"
const PASSWORD_HASH string = "SELECT password FROM USERS WHERE username = ?;"
const CHECK_USER string = "SELECT COUNT(1) FROM USERS WHERE email = ?;"
const INSERT_SESSION_TOKEN = "UPDATE employees SET session_token = ? WHERE email = ?;"

func (conn ForumDB) Insert(query string, values ...interface{}) {
	stmt, err := conn.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err = stmt.Exec(values...); err != nil {
		log.Println(err)
	}
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
