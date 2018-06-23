package main

import (
	"net/http"
	"log"
	"sync"
	"html/template"
	"path/filepath"
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
	t.templ.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
				<head>
					<title>Chat</title>
				</head>
				<body>
					Let's chat!
				</body>
			</html>
		`))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("listenAndServer:", err)
	}
}
