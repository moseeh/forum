package main

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// db, err := sql.Open("sqlite3", "./forum.db")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	server := http.Server{
		Addr:    ":8000",
		Handler: routes(),
	}
	fmt.Printf("Listening on port %s\n", server.Addr)
	server.ListenAndServe()
}
