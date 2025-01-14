package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"forum/internal"
)

type Config struct {
	ClientID     string
	ClientSecret string
}
type GithubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
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

func (app *App) HandleGithubCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	token, err := app.getGithubAccessToken(code)
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}

	githubUser, err := app.getGithubUser(token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	fmt.Println(githubUser.Email)
}

func (app *App) getGithubAccessToken(code string) (string, error) {
	resp, err := http.PostForm("https://github.com/login/oauth/access_token",
		map[string][]string{
			"client_id":     {config.ClientID},
			"client_secret": {config.ClientSecret},
			"code":          {code},
		})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	values := make(map[string]string)
	for _, pair := range strings.Split(string(body), "&") {
		parts := strings.Split(pair, "=")
		values[parts[0]] = parts[1]
	}
	return values["access_token"], nil
}

func (app *App) getGithubUser(token string) (*GithubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GithubUser

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
