package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"forum/internal"
)

type Config struct {
	ClientID     string
	ClientSecret string
}
type GithubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
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

	existingUser, err := app.users.GetUserbYUsername(githubUser.Login)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	user_id := ""
	if existingUser != nil {
		user_id = internal.UUIDGen()
		if existingUser.AuthProvider != "github" {
			githubUser.Login = githubUser.Login + "1"
			err = app.users.InsertUser(user_id, githubUser.Login, "", "", "github", githubUser.AvatarURL)
			if err != nil {
				http.Error(w, "Error adding user", http.StatusInternalServerError)
				return
			}
		} else {
			user_id = existingUser.UserID
		}
	} else {
		user_id = internal.UUIDGen()
		err = app.users.InsertUser(user_id, githubUser.Login, "", "", "github", githubUser.AvatarURL)
		if err != nil {
			http.Error(w, "Error adding user", http.StatusInternalServerError)
			return
		}
	}
	err = app.AddSession(w, r, githubUser.Login, user_id)
	if err != nil {
		app.clearAuthCookies(w)
		http.Error(w, "Error adding session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func (app *App) AddSession(w http.ResponseWriter, r *http.Request, username, user_id string) error {
	session_token := internal.TokenGen(32)
	csrf_token := internal.TokenGen(32)
	expires := time.Now().Local().Add(2 * time.Hour)
	// set session token
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session_token,
		Expires:  expires,
		Path:     "/", // Make sure the path is set to "/"
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	})

	// set csrf token
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrf_token,
		Expires:  expires,
		Path:     "/", // Set the path to "/"
		HttpOnly: false,
	})

	// set username
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Expires:  expires,
		Path:     "/", // Set the path to "/"
		HttpOnly: false,
	})
	// store the tokens in the db
	if err := app.users.NewSession(user_id, session_token, csrf_token, expires.String()); err != nil {
		app.ErrorHandler(w, r, 500)
		return err
	}
	return nil
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
		if strings.Contains(pair, "=") {
			parts := strings.Split(pair, "=")
			values[parts[0]] = parts[1]
		}
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
