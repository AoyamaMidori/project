package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
)

func byeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "bye")
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "hello")
}

func versionHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "<ol>")
	defer fmt.Fprintln(w, "</ol>")
	for k, v := range map[string]interface{}{
		"OS":       runtime.GOOS,
		"ARCH":     runtime.GOARCH,
		"Ver":      runtime.Version(),
		"Compiler": runtime.Compiler,
	} {
		fmt.Fprintf(w, "<li>%s = %v</li>", k, v)
	}
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
	http.HandleFunc("/version", versionHandler)
	http.Serve(ln, nil)
}

func main() {
	fmt.Println("hello")
	makeServer()
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/bye", byeHandler)
	http.ListenAndServe(":50000", nil)
}
