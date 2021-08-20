package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("sleeep")
		time.Sleep(45 * time.Millisecond)
		content := r.URL.Query().Get("test")
		w.Write([]byte(fmt.Sprintf("<h1>Hello 123 from Code %s</h1>", content)))
	})
	http.ListenAndServe(fmt.Sprintf(":%s", "6500"), r)

}
