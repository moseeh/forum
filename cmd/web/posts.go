package main

import (
	"net/http"

	"forum/internal"
)

func (app *App) PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return

	}
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "title and content are required", http.StatusBadRequest)
		return
	}

	user_id, err := app.users.GetUserID(usernameCookie.Value)
	if err != nil {
		http.Error(w, "Error Getting user ID", http.StatusInternalServerError)
		return
	}
	post_id := internal.UUIDGen()
	err = app.users.InsertPost(post_id, title, content, user_id)
	if err != nil {
		http.Error(w, "Error making post", http.StatusInternalServerError)
		return
	}
}
