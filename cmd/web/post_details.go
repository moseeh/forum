package main

import (
	"html/template"
	"net/http"

	db "forum/internal/queries"
)

func (app *App) PostDetailsHandler(w http.ResponseWriter, r *http.Request) {
    postID := r.URL.Query().Get("post_id")
    user_ID := ""
    username := "" // Initialize empty username

    usernameCookie, err := r.Cookie("username")
    if err == nil {
        username = usernameCookie.Value // Set username only if cookie exists
        user_ID, err = app.users.GetUserID(username)
        if err != nil {
            http.Error(w, "Error getting user ID from the database", http.StatusInternalServerError)
            return
        }
    }
    
    post, err := app.users.GetPostDetails(postID, user_ID)
    if err != nil {
        http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
        return
    }
    
    comments, err := app.users.GetPostComments(postID, user_ID)
    if err != nil {
        http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
        return
    }

    // Create a data structure to pass to the template
    data := struct {
        *db.Post
        Comments []db.Comment
        IsLoggedIn bool
        Username string
    }{
        Post:     post,
        Comments: comments,
        IsLoggedIn: user_ID != "",
        Username: username, // Use the initialized username variable
    }
    
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