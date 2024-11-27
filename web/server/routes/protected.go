package routes

import (
	"errors"
	"fmt"
	"forum/web/db"
	"net/http"
)

func authorizer(r *http.Request, dbInstance *db.ForumDB) error {
	// user_email := r.FormValue("email")
	// User exists
	if ok, _ := dbInstance.Exists(db.EMAIL_EXIST, "aaochieng@gmail.com"); !ok {
		return errors.New("Unauthorized 1")
	}

	session_token, err := r.Cookie("session_token")
	if err != nil || session_token.Value == "" {
		return errors.New("Unauthorized 3")
	}
	// csrf := r.Header.Get("X-CSRF-TOKEN")
	// if csrf == "" {
	// 	return errors.New("Unauthorized 4")
	// }
	return nil
}

func protected(w http.ResponseWriter, r *http.Request, dbInstance *db.ForumDB) {
	if r.Method == http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := authorizer(r, dbInstance); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Authorized")
}
