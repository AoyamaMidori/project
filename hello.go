package main

import (
	"fmt"
	"net/http"
	"net"
	"log"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "hello")
}

func listenOnPortAvailable() (net.Listener, string) {
	ln, err := net.Listen("tcp", "localhost:")
	if err != nil {
		log.Fatal(err)
	}
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		log.Fatal(err)
	}
	return ln, port
}

func makeServer() {
	ln, port := listenOnPortAvailable()
	fmt.Println("server is launched on port", port)
	http.HandleFunc("/hello", helloHandler)
	http.Serve(ln, nil)
}

func main() {
	fmt.Println("hello")
	makeServer()
}
