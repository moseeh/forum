package handlers

import "net/http"

func GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get login page"))
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) {

}
