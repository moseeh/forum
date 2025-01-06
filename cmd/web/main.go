package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	internal "forum/internal/queries"

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
	internal.InsertCategories(db)
	server := http.Server{
		Addr:    ":8000",
		Handler: RouteChecker(App.routes()),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	fmt.Printf("Listening on port %s\n", server.Addr)
	select {}
}
