package db

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

type ForumDB struct {
	db *sql.DB
}

func (conn ForumDB) genUUID() string {
	u2, _ := uuid.NewV4()
	return u2.String()
}
