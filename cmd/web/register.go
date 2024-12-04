package main

import (
	"forum/internal"
	"html/template"
	"net/http"
)

func (app *App) register_get(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/templates/register.page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (app *App) register_post(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("./assets/templates/register.page.html"))
	form_errors := map[string][]string{}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm_password := r.FormValue("confirm")

	if err := internal.ValidateUsername(username); err != nil {
		form_errors["username"] = append(form_errors["username"], err.Error())
	}

	if err := internal.ValidateEmail(email); !err {
		form_errors["email"] = append(form_errors["email"], "Invalid email address")
	}

	// validate passwords
	if password == "" || confirm_password == "" {
		form_errors["password"] = append(form_errors["password"], "Empty password not allowed")
	} else if password != confirm_password {
		form_errors["confirm"] = append(form_errors["confirm"], "Passwords do not match")
	} else if len(password) < 8 {
		form_errors["password"] = append(form_errors["password"], "Password must be at least 8 characters")
	}

	if len(form_errors) > 0 {
		tmpl.Execute(w, map[string]interface{}{
			"Errors": form_errors,
		})
	}

	/// check if user exists

	if exists, _ := app.users.UserExists(email); exists {
		http.Error(w, "User already exists", http.StatusFound)
		return
	}

	id := internal.UUIDGen()
	password_hash, _ := internal.HashPassword(password)

	if err := app.users.InsertUser(id, username, email, password_hash); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
