package main

import (
	"net/http"
	"html/template"
	"log"
	"time"
)

var (
	siteUrl = "http://127.0.0.1"
)

type page struct {
	Title string
	Msg string
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")


	t, _ := template.ParseFiles("index.html")
	t.Execute(w, &page{Title: "Convert Image"})

}

func API(w http.ResponseWriter, r *http.Request) {

}

func login(w http.ResponseWriter, r *http.Request) {
	login := r.PostFormValue("login")
	password := r.PostFormValue("pass")

	var user sUser
	user = getUser(login)


	if user.Password == password {
		createToken(login)

		expire := time.Now().AddDate(0, 0, 1)

		cookie := http.Cookie{Name: "token",Value: getToken("acsf"),Expires:expire,Path:"/"}
		http.SetCookie(w,&cookie)
		http.Redirect(w, r, siteUrl,302)

		log.Println("[INFO]: Logged in user:", login)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	createUser(r.PostFormValue("login"),r.PostFormValue("pass"))
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/API/", API)
	http.HandleFunc("/login/", login)
	http.HandleFunc("/register/", register)

	http.ListenAndServe(":80", nil)
}