package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/API/", API)
	http.HandleFunc("/login/", login)
	http.HandleFunc("/register/", register)
	http.HandleFunc("/logout/", logout)

	http.ListenAndServe(":80", nil)
}