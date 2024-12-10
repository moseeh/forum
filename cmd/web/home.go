package main

import (
	"html/template"
	"net/http"

	db "forum/internal/queries"
)

type PageData struct {
	IsLoggedIn bool
	Username   string
	Posts      []db.Post
	Categories []db.Category
}

func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	var data PageData
	data.IsLoggedIn = false

	// Get current user if logged in
	sessionCookie, err := r.Cookie("session_token")
	userID := ""
	if err == nil {
		usernameCookie, err := r.Cookie("username")
		if err == nil {
			if valid, err := app.users.ValidateSession(sessionCookie.Value); valid && err == nil {
				data.IsLoggedIn = true
				data.Username = usernameCookie.Value
				userID, _ = app.users.GetUserID(usernameCookie.Value)
			}
		}
	}

	// Get all categories
	categories, err := app.users.GetAllCategories()
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}
	data.Categories = categories

	// Get posts with categories, likes, and comments
	posts, err := app.users.GetAllPosts(userID) // Pass userID to check if posts are liked by current user
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Posts = posts

	tmpl, err := template.ParseFiles("./assets/templates/index.page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
