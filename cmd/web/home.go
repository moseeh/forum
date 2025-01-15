package main

import (
	"fmt"
	"html/template"
	"net/http"

	db "forum/internal/queries"
)

type PageData struct {
	IsLoggedIn   bool
	Username     string
	Posts        []db.Post
	LikedPosts   []db.Post
	CreatedPosts []db.Post
	Categories   []db.Category
	Trends       []db.CategoryCount
}

func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		IsLoggedIn: false,
		Posts:      []db.Post{},
		Categories: []db.Category{},
		Trends:     []db.CategoryCount{},
	}

	// Get current user if logged in
	sessionCookie, err := r.Cookie("session_token")
	userID := ""
	if err == nil {
		usernameCookie, err := r.Cookie("username")
		if err == nil {
			if valid, err := app.users.ValidateSession(sessionCookie.Value); valid && err == nil {
				data.IsLoggedIn = true
				data.Username = usernameCookie.Value
				userID, _ = app.users.GetUserID(usernameCookie.Value)
			} else {
				fmt.Println(err)
				app.clearAuthCookies(w)
			}
		}
	}

	// Get all categories
	categories, err := app.users.GetAllCategories()
	if err != nil {
		app.ErrorHandler(w, r, 500)
		return
	}
	if categories != nil {
		data.Categories = categories
	}

	// Get posts with categories, likes, and comments
	posts, err := app.users.GetAllPosts(userID)
	if err != nil {
		app.ErrorHandler(w, r, 500)
		return
	}
	if posts != nil {
		data.Posts = posts
	}
	likedPosts, err := app.users.GetLikedPosts(userID)
	if err != nil {
		app.ErrorHandler(w, r, 500)
	}
	if likedPosts != nil {
		data.LikedPosts = likedPosts
	}
	createdPosts, err := app.users.GetCreatedPosts(userID)
	if err != nil {
		app.ErrorHandler(w, r, 500)
	}
	if createdPosts != nil {
		data.CreatedPosts = createdPosts
	}
	trends, err := app.users.TrendingCount()
	if err != nil {
		app.ErrorHandler(w, r, 500)
	}
	if trends != nil {
		Sort(trends)
		if len(trends) > 5 {
			data.Trends = trends[:5]
		} else {
			data.Trends = trends
		}
	}

	tmpl, err := template.ParseFiles("./assets/templates/index.page.html")
	if err != nil {
		app.ErrorHandler(w, r, 500)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		app.ErrorHandler(w, r, 500)
		return
	}
}

func Sort(data []db.CategoryCount) {
	for i := 0; i < len(data)-1; i++ {
		for j := 0; j < len(data)-i-1; j++ {
			if data[j].Count < data[j+1].Count {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}
