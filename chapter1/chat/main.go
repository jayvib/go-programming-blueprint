package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"path/filepath"
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

	t.templ.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./template/chat.html")
		if err != nil {
			log.Println("error:", err)
		}

		err = t.Execute(w, "Jayson")
		if err != nil {
			log.Println("error:", err)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
