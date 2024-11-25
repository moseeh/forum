CREATE TABLE
    IF NOT EXISTS USERS (
        user_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
        username varchar(20) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS POSTS (
        post_id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
        title VARCHAR(255) NOT NULL,
        content VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        author_id VARCHAR(255),
        FOREIGN KEY (author_id) REFERENCES USERS (user_id)
    );

CREATE TRIGGER IF NOT EXISTS update_post_timestamp AFTER
UPDATE ON POSTS FOR EACH ROW BEGIN
UPDATE POSTS
SET
    updated_at = CURRENT_TIMESTAMP
WHERE
    post_id = NEW.post_id;

END;

-- test insert
INSERT INTO
    USERS (user_id, username, email, password)
VALUES
    ('a', 'aaochieng', 'test@email.com', 'encrypted');

-- test insert
INSERT INTO
    POSTS (post_id, title, content, author_id)
VALUES
    (
        '1',
        'First post title',
        'first post content',
        'a'
    );