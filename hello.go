package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/atotto/clipboard"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset-utf8")
	w.WriteHeader(200)
	fmt.Fprintln(w, time.Now())
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset-utf8")
	w.WriteHeader(200)
	fmt.Fprintln(w, r.RemoteAddr)
}

func byeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "bye")
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

func clipboardHandler(w http.ResponseWriter, r *http.Request) {
	data, err := clipboard.ReadAll()
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, err)
		return
	}

	// If query type is specified in query string, Content-Type header
	// is set to the type. Otherwise, Content-Type will be determined
	// by DetectContentType function.
	typ := r.URL.Query().Get("type")
	if typ != "" {
		w.Header().Set("Content-Type", typ)
	}

	fmt.Fprintln(w, data)
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
	http.HandleFunc("/bye", byeHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/ip", ipHandler)

	if !clipboard.Unsupported {
		http.HandleFunc("/cb", clipboardHandler)
	}

	http.Serve(ln, nil)
}

func main() {
	fmt.Println("hello")
	makeServer()
}
