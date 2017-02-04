package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "hello")
}

func main() {
	fmt.Println("hello")
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":50000", nil)
}
