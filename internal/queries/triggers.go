package db

import "log"

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

const COMMENT_LIKES_INCREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS increment_comment_likes
AFTER INSERT ON COMMENT_LIKES
BEGIN
    UPDATE COMMENTS 
    SET likes_count = likes_count + 1
    WHERE comment_id = NEW.comment_id;
END;
`

const LIKES_DECREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS decrement_post_likes
AFTER DELETE ON LIKES
BEGIN
    UPDATE POSTS 
    SET likes_count = likes_count - 1
    WHERE post_id = OLD.post_id;
END;`

const COMMENT_LIKES_DECREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS decrement_comment_likes
AFTER DELETE ON COMMENT_LIKES
BEGIN
    UPDATE COMMENTS 
    SET likes_count = likes_count - 1
    WHERE comment_id = OLD.comment_id;
END;`

const DISLIKES_INCREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS increment_post_dislikes
AFTER INSERT ON DISLIKES
BEGIN
    UPDATE POSTS 
    SET dislikes_count = dislikes_count + 1
    WHERE post_id = NEW.post_id;
END;
`

const COMMENT_DISLIKES_INCREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS increment_comment_dislikes
AFTER INSERT ON COMMENT_DISLIKES
BEGIN
    UPDATE COMMENTS
    SET dislikes_count = dislikes_count + 1
    WHERE comment_id = NEW.comment_id;
END;
`

const DISLIKES_DECREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS decrement_post_dislikes
AFTER DELETE ON DISLIKES
BEGIN
    UPDATE POSTS 
    SET dislikes_count = dislikes_count - 1
    WHERE post_id = OLD.post_id;
END;`

const COMMENT_DISLIKES_DECREMENT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS decrement_comment_dislikes
AFTER DELETE ON COMMENT_DISLIKES
BEGIN
    UPDATE COMMENTS
    SET dislikes_count = dislikes_count - 1
    WHERE comment_id = OLD.comment_id;
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

// Triggers represents all trigger creation statements
var Triggers = []string{
	// Update timestamps
	POST_UPDATE_TRIGGER,
	CATEGORIES_UPDATE_TRIGGER,
	COMMENTS_UPDATE_TRIGGER,

	// Post interactions
	LIKES_INCREMENT_TRIGGER,
	LIKES_DECREMENT_TRIGGER,
	DISLIKES_INCREMENT_TRIGGER,
	DISLIKES_DECREMENT_TRIGGER,
	COMMENTS_INCREMENT_TRIGGER,
	COMMENTS_DECREMENT_TRIGGER,

	// Comment interactions
	COMMENT_LIKES_INCREMENT_TRIGGER,
	COMMENT_LIKES_DECREMENT_TRIGGER,
	COMMENT_DISLIKES_INCREMENT_TRIGGER,
	COMMENT_DISLIKES_DECREMENT_TRIGGER,
}

func (m *UserModel) InitTriggers() {
	for _, trigger := range Triggers {
		if err := m.executeStatement(trigger); err != nil {
			log.Printf("Error creating trigger: %v", err)
		}
	}
}
