package main

import (
	"net/http"

	"forum/internal"
)

func (app *App) CommentsHandler(w http.ResponseWriter, r *http.Request) {
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user_id, err := app.users.GetUserID(usernameCookie.Value)
	if err != nil {
		app.ErrorHandler(w,r,500)
		return
	}
	post_id := r.FormValue("post_id")
	comment := r.FormValue("comment")

	comment_id := internal.UUIDGen()
	err = app.users.InsertComment(comment_id, post_id, user_id, comment)
	if err != nil {
		app.ErrorHandler(w,r,500)
		return
	}
	referer := r.Referer()
	if referer == "" {
		referer = "/home"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}
