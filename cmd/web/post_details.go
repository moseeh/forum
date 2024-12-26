package main

import (
	"html/template"
	"net/http"

	db "forum/internal/queries"
)

func (app *App) PostDetailsHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	user_ID := ""

	usernameCookie, err := r.Cookie("username")
	if err == nil {
		user_ID, err = app.users.GetUserID(usernameCookie.Value)
		if err != nil {
			http.Error(w, "Error getting user ID from the databse", http.StatusInternalServerError)
			return
		}
	}

	post, err := app.users.GetPostDetails(postID, user_ID)
	if err != nil {
		http.Error(w, "DATABSE ERROR", http.StatusInternalServerError)
		return
	}
	comments, err := app.users.GetPostComments(postID)
	if err != nil {
		http.Error(w, "DATABSE ERROR", http.StatusInternalServerError)
		return
	}
	// Create a data structure to pass to the template
	data := struct {
		*db.Post
		Comments []db.Comment
	}{
		Post:     post,
		Comments: comments,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("./assets/templates/post.page.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}
