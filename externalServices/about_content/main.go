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
		time.Sleep(5 * time.Millisecond)
		w.Write([]byte(fmt.Sprintf("<h1>About!!!</h1>")))
	})
	http.ListenAndServe(fmt.Sprintf(":%s", "6500"), r)

}
