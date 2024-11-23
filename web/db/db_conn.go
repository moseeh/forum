package db

import (
	"database/sql"
	"log"
)

const PROJECT_DATABSE = "forum.db"

type ForumDB struct {
	db *sql.DB
}

func NewConnection() *ForumDB {
	return &ForumDB{db: nil}
}

func (conn *ForumDB) db_connection() {
	db, err := sql.Open("sqlite3", PROJECT_DATABSE)
	if err != nil {
		log.Println(err)
		return
	}
	//
	conn.db = db
}
