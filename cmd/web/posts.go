package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"forum/internal"
)

func (app *App) PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.ErrorHandler(w,r, http.StatusMethodNotAllowed)
		return
	}

	// Check authentication
	_, err := r.Cookie("session_token")
	if err != nil {
		app.ErrorHandler(w,r, http.StatusUnauthorized)
		return
	}

	usernameCookie, err := r.Cookie("username")
	if err != nil {
		app.ErrorHandler(w,r, http.StatusUnauthorized)
		return
	}

	// Parse the form to retrieve file and filename
	err = r.ParseMultipartForm(20 << 20) // Max 20MB file size
	if err != nil {
		app.ErrorHandler(w,r, http.StatusBadRequest)
		return
	}

	// Get form values
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categories := r.Form["categories[]"] // Get selected categories

	if title == "" || content == "" {
		app.ErrorHandler(w, r, 400)
		return
	}

	// Get user ID
	userID, err := app.users.GetUserID(usernameCookie.Value)
	if err != nil {
		app.ErrorHandler(w, r, 500)
		return
	}

	postID := internal.UUIDGen()
	var imageName string

	//Handle image upload
	file, header, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		app.ErrorHandler(w, r, 400)
		return
	}

	if file != nil {
		defer file.Close()

		postDir := "assets/static/post_images"

		//add file extension
		ext := filepath.Ext(header.Filename)
		newFilename := postID + ext

		//difine file path
		filepath := filepath.Join(postDir, newFilename)

		// Create the directory if it doesn't exist
		if _, err := os.Stat(postDir); os.IsNotExist(err) {
			os.Mkdir(postDir, os.ModePerm)
		}

		//create the file
		out, err := os.Create(filepath)
		if err != nil {
			app.ErrorHandler(w, r, 500)
			return
		}

		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			app.ErrorHandler(w, r, 500)
			return
		}

		imageName = newFilename
	}

	// Begin transaction
	tx, err := app.users.DB.Begin()
	if err != nil {
		app.ErrorHandler(w, r, 500)
		return
	}
	defer tx.Rollback()

	// Insert post

	err = app.users.InsertPost(tx, postID, title, content, userID, imageName)
	if err != nil {
		app.ErrorHandler(w, r, 500)
		return
	}

	// Insert categories
	for _, categoryID := range categories {
		err = app.users.InsertPostCategory(tx, postID, categoryID)
		if err != nil {
			app.ErrorHandler(w, r, 500)
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		app.ErrorHandler(w, r, 500)
		return
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}
