package main

import (
	"forumed/cmd/handlers"
	"log"
	"net/http"
	"path/filepath"
)

func routes() http.Handler {
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
	mux.HandleFunc("GET /", handlers.HomeHandler)

	//
	mux.HandleFunc("GET /login", handlers.GetLoginHandler)
	mux.HandleFunc("POST /login", handlers.PostLoginHandler)

	//
	return mux
}
