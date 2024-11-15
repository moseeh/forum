package routes

import (
	"fmt"
	"net/http"
	"text/template"
)

func HomeRoute(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println("Error executing template:", err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }
}
