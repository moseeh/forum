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
}

func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.users.GetALlPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data := PageData{
		IsLoggedIn: false,
		Username:   "",
		Posts:      posts,
	}

	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		usernameCookie, err := r.Cookie("username")

		if err == nil {
			if valid, err := app.users.ValidateSession(sessionCookie.Value); valid && err == nil {
				data.IsLoggedIn = true
				data.Username = usernameCookie.Value
			}
		}
	}
	tmpl, err := template.ParseFiles("./assets/templates/index.page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
