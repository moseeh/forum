package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/internal"
	database "forum/internal/queries"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	users *database.UserModel
}

func main() {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = internal.LoadEnvFile(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return
	}

	config = Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}
	googleConfig = GoogleConfig{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}
	App := App{
		users: &database.UserModel{
			DB: db,
		},
	}
	App.users.InitTables()
	database.InsertCategories(db)
	server := http.Server{
		Addr:    ":8000",
		Handler: App.RouteChecker(App.routes()),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
			os.Exit(0)
		}
	}()
	fmt.Printf("Listening on port %s\n", server.Addr)
	select {}
}
