package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CustomSession struct {
	Data map[string]interface{}
}

var globalsession = make(map[string]CustomSession)

func main() {
	http.HandleFunc("/logmein", loginHandler)
	http.Handle("/", http.FileServer(http.Dir("html")))
	http.HandleFunc("/userspace", userspaceHandler)
	http.ListenAndServe(":2304", nil)

}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {
		login := r.FormValue("login")
		password := r.FormValue("password")
		if login == "admin" && password == "admin" {
			newID := uuid.New().String()

			log.Println(newID)
			globalsession[newID] = CustomSession{map[string]interface{}{"login": login, "logintime": time.Now(), "lifetime": time.Now().Add(time.Minute * 10)}}

			http.SetCookie(w, &http.Cookie{
				Name:     "CustomSessionID",
				Value:    newID,
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute * 10),
			})
			http.Redirect(w, r, "/userspace", 301)
			log.Println("User logged in, session created", newID, globalsession[newID])
		} else {
			w.Write([]byte("Wrong login or password"))
		}

	} else {
		w.Write([]byte("POST only"))
	}

}
func userspaceHandler(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header.Get("Cookie")
	if cookie != "" {
		if content, ok := globalsession[cookie[16:]]; ok {
			if content.Data["lifetime"].(time.Time).After(time.Now()) {
				w.Write([]byte("Hello " + content.Data["login"].(string)))
			} else {
				w.Write([]byte("Session expired"))
			}
		} else {
			w.Write([]byte("No session"))
		}
	} else {
		w.Write([]byte("No session"))
	}

}
