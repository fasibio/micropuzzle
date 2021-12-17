package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./public"))
	log.Println("run")
	http.ListenAndServe(fmt.Sprintf(":%s", "6500"), fs)

}
