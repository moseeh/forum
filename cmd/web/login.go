package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"forum/internal"
)

func (app *App) GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		valid, validationErr := app.users.ValidateSession(sessionCookie.Value)
		if validationErr == nil && valid {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
		})
	}
	tmpl, err := template.ParseFiles("./assets/templates/login.page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (app *App) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./assets/templates/login.page.html"))
	form_errors := map[string][]string{}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := internal.ValidateEmail(email); !err {
		form_errors["email"] = append(form_errors["email"], "Invalid email address")
	}

	// validate passwords
	if len(password) < 8 {
		form_errors["password"] = append(form_errors["password"], "Password must be at least 8 characters")
	}

	// confirm if a user exists in the database
	if ok, _ := app.users.UserExists(email); !ok {
		http.Error(w, "User already exists", http.StatusFound)
		return
	}

	// Get password hash
	password_hash, err := app.users.GetPassword(email)
	if err != nil {
		fmt.Println(err)
	}

	// compare password hash
	val := internal.CompareHash(password_hash, password)
	if !val {
		form_errors["password"] = append(form_errors["password"], "Incorrect password")
	}

	if len(form_errors) > 0 {
		tmpl.Execute(w, map[string]interface{}{
			"Errors": form_errors,
		})
		return
	}

	// get username
	user_id, username, err := app.users.GetUsername(email)
	if err != nil {
		fmt.Println(err)
	}

	session_token := internal.TokenGen(32)
	csrf_token := internal.TokenGen(32)
	expires := time.Now().Add(24 * time.Hour)

	// set session token
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session_token,
		Expires:  expires,
		HttpOnly: true,
	})

	// set csfr token
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrf_token,
		Expires:  expires,
		HttpOnly: false,
	})

	// set username
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Expires:  expires,
		HttpOnly: false,
	})

	fmt.Printf("Session token: %s\nCSRF Token: %s\nExpires: %s\n", session_token, csrf_token, expires.Format("2024-11-28 11:13:29"))
	// store the tokens in the db
	if err = app.users.NewSession(user_id, session_token, csrf_token, expires.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("passed")

	http.Redirect(w, r, "/home", http.StatusFound)
}
