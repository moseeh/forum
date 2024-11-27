package routes

import (
	"forum/web/db"
	"net/http"
)

func Routes(w http.ResponseWriter, r *http.Request, dbInstance *db.ForumDB) {
	switch r.URL.Path {
	case "/":
		HomeHandler(w, r)
	case "/register":
		registerhandler(w, r, dbInstance)
	case "/login":
		loginhandler(w, r, dbInstance)
	case "/protected":
		protected(w, r, dbInstance)
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}
