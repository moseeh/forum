package main

import (
	"fmt"
	"forum/web/server/routes"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Define the path to the static files
	staticDir := "./web/static/" // Relative path from `server/main.go` to `web/static`

	// Get the absolute path for better error handling
	absStaticDir, err := filepath.Abs(staticDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path of static directory: %v", err)
	}

	// Serve static files
	fs := http.FileServer(http.Dir(absStaticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Printf("Serving static files from %s on http://localhost:8080/static/", absStaticDir)
	http.HandleFunc("/", routes.RoutesHandler)
	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
