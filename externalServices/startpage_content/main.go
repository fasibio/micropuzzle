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
		time.Sleep(60 * time.Millisecond)
		w.Write([]byte(fmt.Sprintf("<script type=\"module\" src=\"/micro-lib/micropuzzle-components.esm.js\"></script><h1>Startpage!!!</h1>")))
	})
	http.ListenAndServe(fmt.Sprintf(":%s", "6500"), r)

}
