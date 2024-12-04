package db

import "log"

const USERS_TABLE string = `CREATE TABLE
    IF NOT EXISTS USERS (
        user_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
        username varchar(20) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

const POSTS_TABLE = `CREATE TABLE
    IF NOT EXISTS POSTS (
        post_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
        title VARCHAR(255) NOT NULL,
        content VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        author_id VARCHAR(255),
        FOREIGN KEY (author_id) REFERENCES USERS (user_id)
    );`

const POST_UPDATE_TRIGGER = `CREATE TRIGGER IF NOT EXISTS update_post_timestamp AFTER
    UPDATE ON POSTS FOR EACH ROW BEGIN
    UPDATE POSTS
    SET
        updated_at = CURRENT_TIMESTAMP
    WHERE
        post_id = NEW.post_id;

    END;`

const TOKENS = `
    CREATE TABLE IF NOT EXISTS TOKENS(
        id INTEGER PRIMARY KEY ,
        session_token VARCHAR(255) NOT NULL,
        csrf_token VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP NOT NULL,
        user_id VARCHAR(255),
        FOREIGN KEY (user_id) REFERENCES USERS(id)
    );
`

var statements = []string{USERS_TABLE, POSTS_TABLE, POST_UPDATE_TRIGGER, TOKENS}

func (m *UserModel) InitTables() {
	for _, statement := range statements {
		stmt, err := m.DB.Prepare(statement)
		defer stmt.Close()
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := stmt.Exec(); err != nil {
			log.Println(err.Error())
		}
	}
}
