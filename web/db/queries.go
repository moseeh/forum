package db

import (
	"fmt"
	"log"
)

const USER_INSERT string = "INSERT INTO USERS (user_id, username, email, password) VALUES (?,?,?,?);"

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
