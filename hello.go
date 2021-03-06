package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
)

var useCopy = flag.Bool("copy", true, "server url will be copied into the clipbaord.")

var useAutoOpen = flag.Bool("autoopen", false, "launched server will be opened on your web browser")

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

func randomHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Itoa(r.URL.Query().Get("n"))
	if err != nil {
		n = 32
	}
	if n > 1024 {
		n = 1024
	}

	b := make([]byte, n)
	_, err = rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(b)
}

func primeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	q := r.URL.Query()
	bits, err := strconv.Itoa(q.Get("bits"))
	if err != nil {
		bits = 8
	}
	if bits < 2 {
		fmt.Fprintln(w, "bits must be greater than 2")
		return
	}
	if bits > 1024 {
		bits = 1024
	}

	v := interpretBool(q.Get("v"))
	t := time.Now()

	p, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Fprintln(w, p)
}

var interpretedBool = map[string]bool{
	"true": true,
	"yes":  true,
	"on":   true,
	"1":    true,
}

func interpretBool(s string) bool {
	return interpretedBool[s]
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

func startFile(command string, args ...string) error {
	argv := make([]string, len(args)+3)
	argv = append(argv, "/c", "start")
	argv = append(argv, args...)
	argv = append(argv, command)
	return exec.Command("cmd", argv...).Run()
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

	url := "http://localhost:" + port
	if *useCopy && !clipboard.Unsupported {
		err := clipboard.WriteAll(url)
		if err == nil {
			fmt.Println("server url is copied")
		}
	}
	if *useAutoOpen {
		startFile(url)
	}

	http.Serve(ln, nil)
}

func onWin() bool {
	return runtime.GOOS == "windows"
}

func main() {
	fmt.Println("hello")
	makeServer()
}
