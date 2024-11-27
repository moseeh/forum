package routes

import (
	"forum/web/db"
	"forum/web/server/utils"
	"net/http"
	"text/template"
	"time"
)

func loginhandler(w http.ResponseWriter, r *http.Request, dbInstance *db.ForumDB) {
	var tmpl = template.Must(template.ParseFiles("web/templates/auth/login.html"))
	if r.Method == http.MethodGet {
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		//Form errors
		form_errors := map[string][]string{}

		//Parse form values
		if err := r.ParseForm(); err != nil {
			http.Error(w, "400 Bad Request: Unable to parse form", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if err := utils.ValidateEmail(email); err != true {
			form_errors["email"] = append(form_errors["email"], "Invalid email address")
		}

		// validate passwords
		if len(password) < 8 {
			form_errors["password"] = append(form_errors["password"], "Password must be at least 8 characters")
		}

		//confirm if a user exists in the database
		if ok, _ := dbInstance.Exists(db.CHECK_USER, email); !ok {
			form_errors["email"] = append(form_errors["email"], "Incorrect email address")
		}

		//Get password hash
		password_hash, _ := dbInstance.GetStringValue(db.PASSWORD_HASH, email, password)

		// compare password hash
		if !utils.CompareHash(password_hash, password) {
			form_errors["password"] = append(form_errors["password"], "Incorrect password")
		}

		if len(form_errors) > 0 {
			tmpl.Execute(w, map[string]interface{}{
				"Errors": form_errors,
			})
		}

		session_token := utils.TokenGen(32)

		http.SetCookie(w, &http.Cookie{
			Name:    "session token",
			Value:   session_token,
			Expires: time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
