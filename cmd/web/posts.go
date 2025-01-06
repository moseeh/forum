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
		app.ErrorHandler(w,r,400)
		return
	}

	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["categories[]"] // Get selected categories

	if title == "" || content == "" {
		app.ErrorHandler(w,r,400)
		return
	}

	// Get user ID
	userID, err := app.users.GetUserID(usernameCookie.Value)
	if err != nil {
		app.ErrorHandler(w,r,500)
		return
	}

	// Begin transaction
	tx, err := app.users.DB.Begin()
	if err != nil {
		app.ErrorHandler(w,r,500)
		return
	}
	defer tx.Rollback()

	// Insert post
	postID := internal.UUIDGen()
	err = app.users.InsertPost(tx, postID, title, content, userID)
	if err != nil {
		app.ErrorHandler(w,r,500)
		return
	}

	// Insert categories
	for _, categoryID := range categories {
		err = app.users.InsertPostCategory(tx, postID, categoryID)
		if err != nil {
			app.ErrorHandler(w,r,500)
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		app.ErrorHandler(w,r,500)
		return
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}
