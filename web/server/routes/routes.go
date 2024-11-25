package routes

import "net/http"

func Routes(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		HomeHandler(w, r)
	}
}
