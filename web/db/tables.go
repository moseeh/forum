package db

const USERS_TABLE string = `CREATE TABLE
    IF NOT EXISTS USERS (
        user_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
        username varchar(20),
        email VARCHAR(255),
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

var statements = []string{USERS_TABLE, POSTS_TABLE, POST_UPDATE_TRIGGER}
