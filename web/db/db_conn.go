package db

import "database/sql"

const PROJECT_DATABSE = "forum.db"

func connection(db_name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", db_name)
	if err != nil {
		return nil, err
	}
	return db, nil
}
