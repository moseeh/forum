package main

import (
	"fmt"
	"net/http"

	"forum/internal"
)

func (app *App) CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Referer()
	if referer == "" {
		referer = "/home"
	}
	commentID := r.URL.Query().Get("comment_id")
	if commentID == "" {
		http.Error(w, "comment ID is reqired", http.StatusBadRequest)
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
	dislike, err := app.users.UserDislikeOnCommentExists(commentID, user_ID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}
	if dislike {
		err = app.users.DeleteCommentDislike(commentID, user_ID)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}
	}
	exists, err := app.users.UserLikeOnCommentExists(commentID, user_ID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	if exists {
		err = app.users.DeleteCommentLike(commentID, user_ID)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}
	likeID := internal.UUIDGen()

	err = app.users.InsertCommentLike(likeID, commentID, user_ID)
	if err != nil {
		http.Error(w, "Error adding like", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func (app *App) CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Referer()
	if referer == "" {
		referer = "/home"
	}
	commentID := r.URL.Query().Get("comment_id")
	if commentID == "" {
		http.Error(w, "comment ID is reqired", http.StatusBadRequest)
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

	like, err := app.users.UserLikeOnCommentExists(commentID, user_ID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}
	if like {
		err = app.users.DeleteCommentLike(commentID, user_ID)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}
	}
	exists, err := app.users.UserDislikeOnCommentExists(commentID, user_ID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	if exists {
		err = app.users.DeleteCommentDislike(commentID, user_ID)
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}
	dislikeID := internal.UUIDGen()

	err = app.users.InsertCommentDislike(dislikeID, commentID, user_ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error adding dislike", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, referer, http.StatusSeeOther)
}
