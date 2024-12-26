package main

import (
	"net/http"

	"forum/internal"
)

func (app *App) LikesHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Referer()
	if referer == "" {
		referer = "/home"
	}
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
	dislike, err := app.users.UserDislikeOnPostExists(postID, user_ID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}
	if dislike {
		err = app.users.DeleteDislike(postID, user_ID)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}
	}
	exists, err := app.users.UserLikeOnPostExists(postID, user_ID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	if exists {
		err = app.users.DeleteLike(postID, user_ID)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}
	likeID := internal.UUIDGen()

	err = app.users.InsertLike(likeID, postID, user_ID)
	if err != nil {
		http.Error(w, "Error adding like", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, referer, http.StatusSeeOther)
}
