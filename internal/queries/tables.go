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

const COMMENTS_LIKES_TABLE = `CREATE TABLE IF NOT EXISTS COMMENT_LIKES (
    like_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    comment_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (comment_id) REFERENCES COMMENTS (comment_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES USERS (user_id) ON DELETE CASCADE,
    UNIQUE(comment_id, user_id)  -- Prevents multiple likes from same user
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

const COMMENTS_DISLIKES_TABLE = `CREATE TABLE IF NOT EXISTS COMMENT_DISLIKES (
    dislike_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    comment_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (comment_id) REFERENCES COMMENTS (comment_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES USERS (user_id) ON DELETE CASCADE,
    UNIQUE(comment_id, user_id)  -- Prevents multiple dislikes from same user
);
`

const COMMENTS_TABLE = `CREATE TABLE IF NOT EXISTS COMMENTS (
    comment_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    post_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    likes_count INTEGER DEFAULT 0,  
    dislikes_count INTEGER DEFAULT 0,  
    FOREIGN KEY (post_id) REFERENCES POSTS (post_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES USERS (user_id) ON DELETE CASCADE
);`

const TOKENS = `
    CREATE TABLE IF NOT EXISTS TOKENS(
        id INTEGER PRIMARY KEY ,
        session_token VARCHAR(255) NOT NULL,
        csrf_token VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP NOT NULL,
        user_id VARCHAR(255),
        FOREIGN KEY (user_id) REFERENCES USERS(user_id)
    );
`

// Tables represents all table creation statements
var Tables = []string{
	USERS_TABLE,
	CATEGORIES_TABLE,
	POSTS_TABLE,
	POST_CATEGORIES_TABLE,
	COMMENTS_TABLE,
	LIKES_TABLE,
	COMMENTS_LIKES_TABLE,
	DISLIKES_TABLE,
	COMMENTS_DISLIKES_TABLE,
	TOKENS,
}

func (m *UserModel) InitTables() {
	// Create tables first
	for _, statement := range Tables {
		if err := m.executeStatement(statement); err != nil {
			log.Printf("Error creating table: %v", err)
		}
	}
	// Create triggers after tables
	m.InitTriggers()
}

func (m *UserModel) executeStatement(statement string) error {
	stmt, err := m.DB.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return err
}
