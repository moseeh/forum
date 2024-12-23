package main

import (
	"fmt"
	"net/http"

	"forum/internal"
)

func (app *App) DislikesHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Post ID is reqired", http.StatusBadRequest)
		return
	}

	usernameCookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user_ID, err := app.users.GetUserID(usernameCookie.Value)
	if err != nil {
		http.Error(w, "Error getting user ID from the databse", http.StatusInternalServerError)
		return
	}

	like, err := app.users.UserLikeOnPostExists(postID, user_ID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}
	if like {
		err = app.users.DeleteLike(postID, user_ID)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}
	}
	exists, err := app.users.UserDislikeOnPostExists(postID, user_ID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	if exists {
		err = app.users.DeleteDislike(postID, user_ID)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	dislikeID := internal.UUIDGen()

	err = app.users.InsertDislike(dislikeID, postID, user_ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error adding dislike", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
