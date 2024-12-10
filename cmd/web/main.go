package main

import (
	"database/sql"
	"fmt"
	internal "forum/internal/queries"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	users *internal.UserModel
}

func main() {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	App := App{
		users: &internal.UserModel{
			DB: db,
		},
	}

	App.users.InitTables()
	server := http.Server{
		Addr:    ":8000",
		Handler: App.routes(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	fmt.Printf("Listening on port %s\n", server.Addr)
	select {}
}
