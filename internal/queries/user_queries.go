package db

import (
	"database/sql"
	"fmt"
	"log"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) InsertUser(id, username, email, password string) error {
	const USER_INSERT string = "INSERT INTO USERS (user_id, username, email, password) VALUES (?,?,?,?);"
	stmt, err := m.DB.Prepare(USER_INSERT)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(id, username, email, password)
	if err != nil {
		log.Println(res)
		return err
	}
	return nil
}

func (m *UserModel) UserExists(email string) (bool, error) {
	const USER_EXISTS string = "SELECT COUNT(1) FROM USERS WHERE  email = ?;"
	stmt, err := m.DB.Prepare(USER_EXISTS)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *UserModel) GetPassword(email string) (string, error) {
	const PASSWORD_HASH string = "SELECT password FROM USERS WHERE email = ?;"
	stmt, err := m.DB.Prepare(PASSWORD_HASH)
	if err != nil {
		return "", nil
	}
	defer stmt.Close()

	var passwordHash string

	if err = stmt.QueryRow(email).Scan(&passwordHash); err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (m *UserModel) GetUsername(email string) (string, string, error) {
	const USERNAME string = "SELECT user_id,username FROM USERS WHERE email = ?;"
	stmt, err := m.DB.Prepare(USERNAME)
	if err != nil {
		return "", "", nil
	}
	defer stmt.Close()

	var username string
	var user_id string

	if err = stmt.QueryRow(email).Scan(&user_id, &username); err != nil {
		return "", "", err
	}
	return user_id, username, nil
}

func (m *UserModel) GetUserID(username string) (string, error) {
	fmt.Println(username)
	const USERNAME string = "SELECT user_id FROM USERS WHERE username = ?;"
	stmt, err := m.DB.Prepare(USERNAME)
	if err != nil {
		return "", nil
	}
	defer stmt.Close()

	var user_id string

	if err = stmt.QueryRow(username).Scan(&user_id); err != nil {
		return "", err
	}
	return user_id, nil
}
