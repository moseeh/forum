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
	data := PageData{
		IsLoggedIn: false,
		Posts:      []db.Post{},
		Categories: []db.Category{},
	}

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
	if categories != nil {
		data.Categories = categories
	}

	// Get posts with categories, likes, and comments
	posts, err := app.users.GetAllPosts(userID) // Pass userID to check if posts are liked by current user
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if posts != nil {
		data.Posts = posts
	}

	tmpl, err := template.ParseFiles("./assets/templates/index.page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
