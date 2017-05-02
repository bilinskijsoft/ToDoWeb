package main

import (
	"net/http"
)

func main() {
	finish := make(chan bool)

	server443 := http.NewServeMux()

	server443.HandleFunc("/", index)
	server443.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	server443.HandleFunc("/API/", API)
	server443.HandleFunc("/login/", login)
	server443.HandleFunc("/register/", register)
	server443.HandleFunc("/logout/", logout)

	go http.ListenAndServeTLS(":443", "server.crt", "server.key", server443)

	server80 := http.NewServeMux()
	server80.HandleFunc("/", redirectToHttps)

	go http.ListenAndServe(":80", server80)

	<-finish
}
