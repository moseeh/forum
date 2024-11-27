package routes

import (
	"forum/web/db"
	"forum/web/server/utils"
	"net/http"
	"text/template"
)

func registerhandler(w http.ResponseWriter, r *http.Request, dbInstance *db.ForumDB) {
	var tmpl = template.Must(template.ParseFiles("web/templates/auth/register.html"))
	if r.Method == http.MethodPost {
		//Form errors
		form_errors := map[string][]string{}

		//Parse form values
		if err := r.ParseForm(); err != nil {
			http.Error(w, "400 Bad Request: Unable to parse form", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirm_password := r.FormValue("confirm")

		if err := utils.ValidateUsername(username); err != nil {
			form_errors["username"] = append(form_errors["username"], err.Error())
		}

		if err := utils.ValidateEmail(email); err != true {
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

		//check if the user already exists in the database
		ok, err := dbInstance.Exists(db.USER_EXISTS, username, password)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if ok {
			http.Error(w, "Record already exists!", http.StatusConflict)
			return
		}

		if len(form_errors) > 0 {
			tmpl.Execute(w, map[string]interface{}{
				"Errors": form_errors,
			})
		}

		//Insert user to the database
		new_uuid := utils.UUIDGen()
		hashed_password, _ := utils.HashPassword(password)
		dbInstance.Insert(db.USER_INSERT, new_uuid, username, email, hashed_password)

	} else if r.Method == http.MethodGet {
		tmpl.Execute(w, nil)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
