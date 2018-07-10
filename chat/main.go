package main

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	gomniauth.SetSecurityKey("AUTHKEY")
	gomniauth.WithProviders(
		github.New("ce78b11d453f2028f5c1", "156329102b712c05d5fd5a401e89876b8cdcf4e4", "http://localhost:3000/auth/callback/github"),
	)
	var addr = flag.String("addr", ":8080", "The address of the application ")
	flag.Parse()
	r := newRoom(UseGravatar)
	http.Handle("/chat", mustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/auth/", loginHandler)
    http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
    	http.SetCookie(w, &http.Cookie{
    		Name: "auth",
    		Value: "",
    		Path: "/",
    		MaxAge: -1,
		})
    	w.Header().Set("Location", "/chat")
    	w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("listenAndServer:", err)
	}
}
