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
		app.clearAuthCookies(w)
	}
	tmpl, err := template.ParseFiles("./assets/templates/login.page.html")
	if err != nil {
		app.ErrorHandler(w, r, 500)
		return
	}
	tmpl.Execute(w, nil)
}

func (app *App) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./assets/templates/login.page.html"))
	form_errors := map[string][]string{}
	if err := r.ParseForm(); err != nil {
		app.ErrorHandler(w, r, 400)
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
		form_errors["email"] = append(form_errors["email"], "User does not exist")
	} else {
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
	expires := time.Now().Local().Add(2 * time.Hour)

	// set session token
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session_token,
		Expires:  expires,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	})

	// set csrf token
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

	// store the tokens in the db
	if err = app.users.NewSession(user_id, session_token, csrf_token, expires.String()); err != nil {
		app.ErrorHandler(w, r, 500)
		return
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}

func (app *App) clearAuthCookies(w http.ResponseWriter) {
	cookies := []string{"session_token", "csrf_token", "username"}
	for _, cookieName := range cookies {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}
}
