package main

import (
	"log"
	"net/http"
	"time"
)

func (app *App) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session cookie
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		// Delete the session from the database
		_, err = app.users.DB.Exec(`DELETE FROM TOKENS WHERE session_token = ?`, sessionCookie.Value)
		if err != nil {
			log.Printf("Error deleting session: %v", err)
		}
	}

	// Clear all cookies
	cookies := []string{"session_token", "csrf_token", "username"}
	for _, cookieName := range cookies {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-24 * time.Hour),
			HttpOnly: true,
		})
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
