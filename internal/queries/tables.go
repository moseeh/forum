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

const CATEGORIES_TABLE string = `CREATE TABLE IF NOT EXISTS CATEGORIES (
    category_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const POSTS_TABLE = `CREATE TABLE IF NOT EXISTS POSTS (
    post_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    author_id VARCHAR(255),
    likes_count INTEGER DEFAULT 0,  -- Counter for quick access to total likes
    dislikes_count INTEGER DEFAULT 0,  -- Counter for quick access to total dislikes
    comments_count INTEGER DEFAULT 0,  -- Counter for quick access to total comments
    FOREIGN KEY (author_id) REFERENCES USERS (user_id)
);`

const POST_CATEGORIES_TABLE = `CREATE TABLE IF NOT EXISTS POST_CATEGORIES (
    post_id VARCHAR(255) NOT NULL,
    category_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES POSTS (post_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES CATEGORIES (category_id) ON DELETE CASCADE
);`

const LIKES_TABLE = `CREATE TABLE IF NOT EXISTS LIKES (
    like_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    post_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES POSTS (post_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES USERS (user_id) ON DELETE CASCADE,
    UNIQUE(post_id, user_id)  -- Prevents multiple likes from same user
);
`

const DISLIKES_TABLE = `CREATE TABLE IF NOT EXISTS DISLIKES (
    dislike_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    post_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES POSTS (post_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES USERS (user_id) ON DELETE CASCADE,
    UNIQUE(post_id, user_id)  -- Prevents multiple dislikes from same user
);
`

const COMMENTS_TABLE = `CREATE TABLE IF NOT EXISTS COMMENTS (
    comment_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    post_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    parent_comment_id VARCHAR(255),  -- For nested comments, NULL if top-level
    FOREIGN KEY (post_id) REFERENCES POSTS (post_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES USERS (user_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES COMMENTS (comment_id) ON DELETE CASCADE
);`

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

const POST_UPDATE_TRIGGER = `CREATE TRIGGER IF NOT EXISTS update_posts_updated_at
AFTER UPDATE ON POSTS
BEGIN
    UPDATE POSTS SET updated_at = CURRENT_TIMESTAMP
    WHERE post_id = NEW.post_id;
END;`

const CATEGORIES_UPDATE_TRIGGER = `CREATE TRIGGER IF NOT EXISTS update_categories_updated_at
AFTER UPDATE ON CATEGORIES
BEGIN
    UPDATE CATEGORIES SET updated_at = CURRENT_TIMESTAMP
    WHERE category_id = NEW.category_id;
END;
`

const COMMENTS_UPDATE_TRIGGER = `CREATE TRIGGER IF NOT EXISTS update_comments_updated_at
AFTER UPDATE ON COMMENTS
BEGIN
    UPDATE COMMENTS SET updated_at = CURRENT_TIMESTAMP
    WHERE comment_id = NEW.comment_id;
END;`

const LIKES_INCREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS increment_post_likes
AFTER INSERT ON LIKES
BEGIN
    UPDATE POSTS 
    SET likes_count = likes_count + 1
    WHERE post_id = NEW.post_id;
END;
`

const LIKES_DECREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS decrement_post_likes
AFTER DELETE ON LIKES
BEGIN
    UPDATE POSTS 
    SET likes_count = likes_count - 1
    WHERE post_id = OLD.post_id;
END;`

const DISLIKES_INCREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS increment_post_dislikes
AFTER INSERT ON DISLIKES
BEGIN
    UPDATE POSTS 
    SET dislikes_count = dislikes_count + 1
    WHERE post_id = NEW.post_id;
END;
`

const DISLIKES_DECREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS decrement_post_dislikes
AFTER DELETE ON DISLIKES
BEGIN
    UPDATE POSTS 
    SET dislikes_count = dislikes_count - 1
    WHERE post_id = OLD.post_id;
END;`

const COMMENTS_INCREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS increment_post_comments
AFTER INSERT ON COMMENTS
BEGIN
    UPDATE POSTS 
    SET comments_count = comments_count + 1
    WHERE post_id = NEW.post_id;
END;`

const COMMENTS_DECREMENT_TRIGGER = `
CREATE TRIGGER IF NOT EXISTS decrement_post_comments
    AFTER DELETE ON COMMENTS
    BEGIN
        UPDATE POSTS 
        SET comments_count = comments_count - 1
        WHERE post_id = OLD.post_id;
    END;`

var statements = []string{USERS_TABLE, POSTS_TABLE, CATEGORIES_TABLE, POST_CATEGORIES_TABLE, LIKES_TABLE, DISLIKES_TABLE, DISLIKES_INCREMENT_TRIGGER, DISLIKES_DECREMENT_TRIGGER, COMMENTS_TABLE, POST_UPDATE_TRIGGER, CATEGORIES_UPDATE_TRIGGER, COMMENTS_UPDATE_TRIGGER, LIKES_INCREMENT_TRIGGER, LIKES_DECREMENT_TRIGGER, COMMENTS_INCREMENT_TRIGGER, COMMENTS_DECREMENT_TRIGGER, TOKENS}

func (m *UserModel) InitTables() {
	for _, statement := range statements {
		stmt, err := m.DB.Prepare(statement)
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := stmt.Exec(); err != nil {
			log.Println(err.Error())
		}
		stmt.Close()
	}
}
