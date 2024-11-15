package main

import (
	"fmt"
	"forum/web/server/routes"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", routes.HomeRoute)
	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
