package db

import (
	"fmt"
)

func (m *UserModel) InsertPost(post_id, title, content, user_id string) error {
	const POST_INSERT string = "INSERT INTO POSTS (post_id, title, content, author_id) VALUES (?,?,?,?)"
	stmt, err := m.DB.Prepare(POST_INSERT)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(post_id, title, content, user_id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
