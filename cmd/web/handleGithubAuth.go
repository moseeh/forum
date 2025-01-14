package main

import (
	"fmt"
	"net/http"

	"forum/internal"
)

type Config struct {
	ClientID     string
	ClientSecret string
}

var config Config

func (app *App) HandleGithubAuth(w http.ResponseWriter, r *http.Request) {

	state := internal.TokenGen(10)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   3600,
		HttpOnly: true,
	})

	// redirect to github
	githubURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s",
		config.ClientID,
		"http://localhost:8000/auth/github/callback",
		state,
	)
	http.Redirect(w, r, githubURL, http.StatusTemporaryRedirect)
}
