package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var allowedRoutes = map[string]bool{
	"/":                true,
	"/home":            true,
	"/login":            true,
	"/register":        true,
	"/logout":          true,
	"/post/like":       true,
	"/post/dislike":    true,
	"/post/details":    true,
	"/posts/create":     true,
	"/comment":         true,
	"/comment/like":    true,
	"/comment/dislike": true,
}

// RouteChecker is a middleware that checkes allowed routes
func (app *App)RouteChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			// Static(w,r)
			next.ServeHTTP(w, r)
			return
		}

		if _, ok := allowedRoutes[r.URL.Path]; !ok {
			app.ErrorHandler(w,r,404)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *App) routes() http.Handler {
	//
	staticDir := "./assets/static/"
	absStaticDir, err := filepath.Abs(staticDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path of static directory: %v", err)
	}
	fs := http.FileServer(http.Dir(absStaticDir))
	//

	mux := http.NewServeMux()

	///
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /", app.HomeHandler)
	mux.HandleFunc("GET /home", app.HomeHandler)

	//
	mux.HandleFunc("GET /login", app.GetLoginHandler)
	mux.HandleFunc("POST /login", app.PostLoginHandler)

	//
	mux.HandleFunc("GET /register", app.register_get)
	mux.HandleFunc("POST /register", app.register_post)

	//
	mux.HandleFunc("GET /logout", app.LogoutHandler)
	mux.HandleFunc("GET /post/like", app.LikesHandler)
	mux.HandleFunc("GET /post/dislike", app.DislikesHandler)

	//
	mux.HandleFunc("GET /post/details", app.PostDetailsHandler)
	mux.HandleFunc("POST /posts/create", app.PostsHandler)
	mux.HandleFunc("POST /comment", app.CommentsHandler)

	//
	mux.HandleFunc("GET /comment/like", app.CommentLikeHandler)
	mux.HandleFunc("GET /comment/dislike", app.CommentDislikeHandler)

	return mux
}
