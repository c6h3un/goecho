package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are from %s!", r.RemoteAddr)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8888", nil)
}