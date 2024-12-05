package main

import (
	"fmt"
	"net/http"
)

func (app *App) PostsHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		usernameCookie, err := r.Cookie("username")
		if err == nil {
			fmt.Println(usernameCookie.Value)
			fmt.Println(sessionCookie.Value)
			title := r.FormValue("title")
			content := r.FormValue("content")
			fmt.Println(title, content)
		}

	}
}
