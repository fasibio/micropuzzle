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
		time.Sleep(5 * time.Second)
		content := r.URL.Query().Get("test")
		w.Write([]byte("Hello from Code " + content))
	})
	http.ListenAndServe(fmt.Sprintf(":%s", "6500"), r)

}
