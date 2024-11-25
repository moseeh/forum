package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"forum/web/db"
	"forum/web/server/routes"
)

func main() {
	// Databases initialization
	// db instance

	database := db.NewConnection()
	database.Connect()
	database.InitTables()
	database.Insert(db.USER_INSERT, "c", "aaochieng", "test@gmail.com", "encrypted")

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
	http.HandleFunc("/", routes.Routes)
	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
