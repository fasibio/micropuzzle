package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

//go:embed index.html
var htmlContent string

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("sleeep")
		time.Sleep(60 * time.Millisecond)
		w.Write([]byte(fmt.Sprintf(htmlContent)))
	})
	http.ListenAndServe(fmt.Sprintf(":%s", "6500"), r)

}
