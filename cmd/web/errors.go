package main

import (
	"html/template"
	"net/http"
)

var ErrorPages = map[int][]string{
	404: {"404 Not Found", "The page you are looking for could not be found."},
	405: {"405 Method not allowed", "The method is not allowed"},
	401: {"401 Unauthorized", "You are not authorized"},
	400: {"400 Bad Request", "Your request cannot be processed due to malformed syntax."},
	500: {"500 Internal Server Error", "The server encountered an unexpected condition."},
	409: {"409 Conflict", "Username already taken"},
}

func (app *App) ErrorHandler(w http.ResponseWriter, r *http.Request, Status int) {
	user_ID := ""
	username := "" // Initialize empty username
	status := ""
	message := ""

	usernameCookie, err := r.Cookie("username")
	if err == nil {
		username = usernameCookie.Value // Set username only if cookie exists
		user_ID, err = app.users.GetUserID(username)
		if err != nil {
			http.Error(w, "Error getting user ID from the database", http.StatusInternalServerError)
			return
		}
	}

	if val, ok := ErrorPages[Status]; ok {
		status = val[0]
		message = val[1]
	}

	// Create a data structure to pass to the template
	data := struct {
		DisplayStatus string
		Message       string
		IsLoggedIn    bool
		Username      string
	}{
		DisplayStatus: status,
		Message:       message,
		IsLoggedIn:    user_ID != "",
		Username:      username, // Use the initialized username variable
	}

	tmpl, err := template.ParseFiles("./assets/templates/error.page.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(Status)
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}
