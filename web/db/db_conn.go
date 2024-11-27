package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const PROJECT_DATABSE = "./forum.db"

func NewConnection() *ForumDB {
	return &ForumDB{db: nil}
}

func (conn *ForumDB) Connect() {
	db, err := sql.Open("sqlite3", PROJECT_DATABSE)
	if err != nil {
		log.Println(err)
		return
	}
	//
	conn.db = db
}

func (conn ForumDB) InitTables() {
	for _, statement := range statements {
		stmt, err := conn.db.Prepare(statement)
		defer stmt.Close()
		if err != nil {
			log.Println(err)
		}
		if _, err := stmt.Exec(); err != nil {
			log.Println(err.Error())
		}
	}
}
