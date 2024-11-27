package routes

import (
	"fmt"
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

		if err := utils.ValidateEmail(email); !err {
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
		password_hash, err := dbInstance.GetStringValue(db.PASSWORD_HASH, email)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(password_hash)

		// compare password hash
		val := utils.CompareHash(password_hash, password)
		if !val {
			form_errors["password"] = append(form_errors["password"], "Incorrect password")
		}

		if len(form_errors) > 0 {
			tmpl.Execute(w, map[string]interface{}{
				"Errors": form_errors,
			})
			return
		}

		//get user id
		user_id, err := dbInstance.GetStringValue(db.USER_ID, email)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("User id : %s\n", user_id)
		session_token := utils.TokenGen(32)
		csrf_token := utils.TokenGen(32)
		expires := time.Now().Add(24 * time.Hour)

		// set session token
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    session_token,
			Expires:  expires,
			HttpOnly: true,
		})

		//set csfr token
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrf_token,
			Expires:  expires,
			HttpOnly: false,
		})
		fmt.Printf("Session token: %s\nCSRF Token: %s\nExpires: %s\n", session_token, csrf_token, expires.Format("2024-11-28 11:13:29"))
		// store the tokens in the db
		err = dbInstance.Insert(db.INSERT_TOKENS, user_id, session_token, csrf_token, expires.String())
		http.Redirect(w, r, "/", http.StatusOK)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
