package main

import (
	"fmt"
	"net/http"
)


func main() {
	http.HandleFunc("/", Hello)
	http.ListenAndServe(":8000", nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1> Hello world!</h1>")
}

