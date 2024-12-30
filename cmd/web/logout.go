package main

import (
	"log"
	"net/http"
	"time"
)

func (app *App) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		err = app.users.DeleteSession(sessionCookie.Value)
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
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure: true,
			SameSite: http.SameSiteStrictMode,
		})
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
