package main

import (
	"net/http"
	"log"
	"sync"
	"html/template"
	"path/filepath"
	"flag"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth"
)

type templateHandler struct {
	once		sync.Once
	filename	string
	templ		*template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	gomniauth.SetSecurityKey("AUTHKEY")
	gomniauth.WithProviders(
		github.New("d82f1aa0f24805644133", "62055e500525ae012665f92bc39b2c8444ef4bb3", "http://localhost:3000/auth/callback/github"),
	)
	var addr = flag.String("addr", ":8080", "The address of the application ")
	flag.Parse()
	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("listenAndServer:", err)
	}
}
