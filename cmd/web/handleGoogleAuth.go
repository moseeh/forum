package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forum/internal"
)

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
}

type GoogleUser struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type GoogleTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token"`
}

var googleConfig GoogleConfig

func (app *App) HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	state := internal.TokenGen(10)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   3600,
		HttpOnly: true,
	})

	googleURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&state=%s&response_type=code&scope=openid profile email",
		googleConfig.ClientID,
		"http://localhost:8000/auth/google/callback",
		state,
	)

	http.Redirect(w, r, googleURL, http.StatusTemporaryRedirect)
}

func (app *App) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")

	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := app.getGoogleAccessToken(code)
	if err != nil {
		http.Error(w, "Failed to get Token", http.StatusInternalServerError)
		return
	}

	googleUser, err := app.getGoogleUser(token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	existingUser, err := app.users.GetUserbYUsername(googleUser.Name)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	user_id := ""
	if existingUser != nil {
		if existingUser.AuthProvider != "google" {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		} else {
			user_id = existingUser.UserID
		}
	} else {
		user_id = internal.UUIDGen()
		err = app.users.InsertUser(user_id, googleUser.Name, googleUser.Email, "", "google", googleUser.Picture)
		if err != nil {
			http.Error(w, "Error Adding User", http.StatusInternalServerError)
			return
		}
	}
	err = app.AddSession(w, r, googleUser.Name, user_id)
	if err != nil {
		app.clearAuthCookies(w)
		http.Error(w, "Error adding session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func (app *App) getGoogleAccessToken(code string) (string, error) {
	resp, err := http.PostForm("https://oauth2.googleapis.com/token",
		map[string][]string{
			"client_id":     {googleConfig.ClientID},
			"client_secret": {googleConfig.ClientSecret},
			"code":          {code},
			"redirect_uri":  {"http://localhost:8000/auth/google/callback"},
			"grant_type":    {"authorization_code"},
		})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp GoogleTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", err
	}
	return tokenResp.AccessToken, nil
}

func (app *App) getGoogleUser(token string) (*GoogleUser, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
