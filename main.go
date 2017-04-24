package main

import (
	"net/http"
	"html/template"
	"log"
	"time"
	"fmt"
	"encoding/json"
)

var (
	siteUrl = "http://127.0.0.1"
)

type page struct {
	Title string
	Msg string
}

type apiUserResponse struct {
	User string
	Id int
}

type apiToDoResponse struct {
	user string
	text string
	status int
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")


	t, _ := template.ParseFiles("index.html")
	t.Execute(w, &page{Title: "Convert Image"})

}

func API(w http.ResponseWriter, r *http.Request) {
	if (r.PostFormValue("method")=="getUser") {
		cookie, err := r.Cookie("token")

		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "{}")
			return
		}

		token := cookie.Value
		if (getUserNameByToken(cookie.Value) != "") {
			user := getUser(getUserNameByToken(token))

			log.Println(getUserNameByToken(token))

			response := apiUserResponse{getUserNameByToken(token), user.Id}

			json, err := json.Marshal(response)

			if err != nil {
				log.Println(err)
			}

			fmt.Fprintf(w, string(json))
		}
	} else {
		fmt.Fprintf(w, "{}")
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	login := r.PostFormValue("login")
	password := r.PostFormValue("pass")

	var user sUser
	user = getUser(login)


	if user.Password == password {
		createToken(login)
		token := getToken(login)

		expire := time.Now().AddDate(0, 0, 1)
		cookie := http.Cookie{Name: "token",Value: token ,Expires:expire,Path:"/"}

		http.SetCookie(w,&cookie)
		http.Redirect(w, r, siteUrl,302)

		log.Println("[INFO]: Logged in user:", login)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	createUser(r.PostFormValue("login"),r.PostFormValue("pass"))
}

func logout(w http.ResponseWriter, r *http.Request) {
	expire := time.Now().Add(1 * time.Second)
	cookie := http.Cookie{Name: "token",Value: "" ,Expires:expire,Path:"/"}

	http.SetCookie(w,&cookie)
	http.Redirect(w, r, siteUrl,302)
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/API/", API)
	http.HandleFunc("/login/", login)
	http.HandleFunc("/register/", register)
	http.HandleFunc("/logout/", logout)

	http.ListenAndServe(":80", nil)
}