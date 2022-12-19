package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!!! New changes")
}

func main() {
	mux := http.NewServeMux()

	hh := http.HandlerFunc(helloHandler)
	mux.Handle("/hello", hh)

	p := os.Getenv("HTTP_PORT")
	if p == "" {
		p = "8080"
	}

	log.Printf("Listening at %s\n", p)
	http.ListenAndServe(":" + p, mux)
}