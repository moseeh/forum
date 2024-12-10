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

	// Check authentication
	_, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	usernameCookie, err := r.Cookie("username")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["categories[]"] // Get selected categories

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Get user ID
	userID, err := app.users.GetUserID(usernameCookie.Value)
	if err != nil {
		http.Error(w, "Error getting user ID", http.StatusInternalServerError)
		return
	}

	// Begin transaction
	tx, err := app.users.DB.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert post
	postID := internal.UUIDGen()
	err = app.users.InsertPost(tx, postID, title, content, userID)
	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	// Insert categories
	for _, categoryID := range categories {
		err = app.users.InsertPostCategory(tx, postID, categoryID)
		if err != nil {
			http.Error(w, "Error linking categories", http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Error completing post creation", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}
