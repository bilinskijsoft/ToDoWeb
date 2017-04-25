package main

import (
	"time"
	"net/http"
	"html/template"
	"log"
	"fmt"
	"encoding/json"
	"strings"
	"strconv"
)

type page struct {
	Title string
	Msg   string
}

type apiUserResponse struct {
	User string
	Id   int
}

type apiNotifyStruct struct {
	Title  string
	Text   string
	Status string
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, &page{Title: "Convert Image"})

}

func API(w http.ResponseWriter, r *http.Request) {
	switch method := r.PostFormValue("method"); method {
	case "getUser":
		cookie, err := r.Cookie("token")

		if err != nil {
			fmt.Fprintf(w, "{}")
			return
		}

		token := cookie.Value
		if getUserNameByToken(cookie.Value) != "" {
			user := getUser(getUserNameByToken(token))

			response := apiUserResponse{getUserNameByToken(token), user.Id}

			json, err := json.Marshal(response)

			if err != nil {
				log.Println(err)
			}

			fmt.Fprintf(w, string(json))
		}
	case "getToDoS":
		cookie, err := r.Cookie("token")

		if err != nil {
			fmt.Fprintf(w, "{}")
			return
		}

		token := cookie.Value
		if getUserNameByToken(cookie.Value) != "" {

			toDoS := getToDoS(getUserNameByToken(token))

			json, err := json.MarshalIndent(toDoS, "todo_", "todo")

			if err != nil {
				log.Println(err)
			}

			result := strings.Replace(string(json), "\\", "", -1)
			result = strings.Trim(result, "\"")
			fmt.Fprintf(w, result)
		}
	case "addToDo":
		cookie, err := r.Cookie("token")

		if err != nil {
			fmt.Fprintf(w, "{}")
			return
		}

		token := cookie.Value

		text := r.PostFormValue("text")
		user := getUserNameByToken(token)

		addToDo(user, text)

		log.Println("[INFO]: Added new todo. User:", user)

		if r.PostFormValue("redirect") == "1" {
			http.Redirect(w, r, "http://"+r.Host, 302)
		}
	case "getNotifys":
		_, err := r.Cookie("token")

		if err != nil {
			response := apiNotifyStruct{"Привет!",
						    "Для продолжения, пожалуйста войдите под своим логином, или зарегистрируйтесь!",
						    "success"}

			json, _ := json.Marshal(response)

			fmt.Fprintf(w, string(json))
			return
		}
		fmt.Fprintf(w, "{}")

	case "editToDo":
		_, err := r.Cookie("token")

		if err != nil {
			fmt.Fprintf(w, "{}")
			return
		}

		id, _ := strconv.Atoi(r.PostFormValue("id"))
		text := r.PostFormValue("text")
		status, _ := strconv.Atoi(r.PostFormValue("status"))
		editToDo(id, text, status)
		http.Redirect(w, r, "http://"+r.Host, 302)

	case "getToDoById":
		_, err := r.Cookie("token")

		if err != nil {
			fmt.Fprintf(w, "{}")
			return
		}

		id, _ := strconv.Atoi(r.PostFormValue("id"))

		fmt.Fprintf(w, getToDoById(id))

	case "deleteToDoById":
		_, err := r.Cookie("token")

		if err != nil {
			fmt.Fprintf(w, "{}")
			return
		}

		id, _ := strconv.Atoi(r.PostFormValue("id"))

		deleteToDo(id)
		http.Redirect(w, r, "http://"+r.Host, 302)

	default:
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
		cookie := http.Cookie{Name: "token", Value: token, Expires: expire, Path: "/"}

		http.SetCookie(w, &cookie)

		log.Println("[INFO]: Logged in user:", login)
	} else {
		fmt.Fprintf(w, "Неправильный логин или пароль")
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	login := r.PostFormValue("login")
	password := r.PostFormValue("pass")

	result := createUser(login, password)

	if result != 1 {
		fmt.Fprintf(w, "Пользователь уже существует!")
		return
	}

	createToken(login)
	token := getToken(login)

	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{Name: "token", Value: token, Expires: expire, Path: "/"}

	http.SetCookie(w, &cookie)

	log.Println("[INFO]: Registered user:", login)
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")

	if err != nil {
		fmt.Fprintf(w, "{}")
		return
	}

	log.Println("[INFO]: Logged out user:", getUserNameByToken(cookie.Value))

	expire := time.Now().Add(1 * time.Millisecond)
	cookieUnset := http.Cookie{Name: "token", Value: "", Expires: expire, Path: "/"}

	http.SetCookie(w, &cookieUnset)
	http.Redirect(w, r, "http://"+r.Host, 302)
}
