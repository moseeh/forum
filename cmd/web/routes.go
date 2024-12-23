package main

import (
	"log"
	"net/http"
	"path/filepath"
)

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
	mux.HandleFunc("GET /likes", app.LikesHandler)
	mux.HandleFunc("GET /dislike", app.DislikesHandler)


	//
	mux.HandleFunc("POST /posts/create", app.PostsHandler)

	return mux
}
