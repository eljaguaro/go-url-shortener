package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thanhpk/randstr"
)

var urlmap = make(map[string]string)

func makeshortFunc(long string) string {
	// return "fewe"
	short := randstr.String(8)
	urlmap[short] = long
	return short
}

func geturlFunc(url string) string {
	// return "fewef"
	long, ok := urlmap[url]
	if !ok {
		return "The short url not found"
	}
	return long
}

func geturlHandle(rw http.ResponseWriter, r *http.Request) {
	// if r.Host != "localhost:8080" {
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	long := geturlFunc(chi.URLParam(r, "id"))
	if long == "The short url not found" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	// rw.WriteHeader(http.StatusTemporaryRedirect)
	// rw.Header().Set("Location", long)
	// http.Redirect(rw, r, long, http.StatusPermanentRedirect)
	http.Redirect(rw, r, long, 307)
	// rw.Write([]byte("few"))
}

type URL struct {
	URL string
}

func makeshortHandle(rw http.ResponseWriter, r *http.Request) {
	// if r.Host != "localhost:8080" {
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	var url URL
	b, err := io.ReadAll(r.Body)
	url.URL = string(b)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}
	// http.RedirectHandler()
	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "text/plain")
	rw.Write([]byte("http://localhost:8080/" + makeshortFunc(url.URL)))
}

// func rhandle(rw http.ResponseWriter, r *http.Request) {
// 	http.Redirect(rw, r, "http://abc.ru/", 309)

// }

// curl -X POST 'http://localhost:8080/' -H "text/plain" -d '{"URL": "abc"}'

func main() {
	// urlmap := make(map[string]string)
	r := chi.NewRouter()
	r.Post("/", makeshortHandle)
	r.Get("/{id}", geturlHandle)
	// r.Get("/fw", rhandle)
	run := flag.String("a", "localhost:8080", "адрес запуска http-сервера")
	surl := flag.String("b", "http://localhost:8000/qsd54gFg", "безовый адрес результирующего URL")
	flag.Parse()
	fmt.Println("address to run the server:", *run)
	fmt.Println("server address and shorturl", *surl)
	log.Fatal(http.ListenAndServe(":8080", r))
}

// package main

// import (
// 	"encoding/json"
// 	"io"
// 	"log"
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/thanhpk/randstr"
// )

// var urlmap = make(map[string]string)

// func makeshortFunc(long string) string {
// 	// return "fewe"
// 	short := randstr.String(8)
// 	urlmap[short] = long
// 	return short
// }

// func geturlFunc(url string) string {
// 	// return "fewef"
// 	long, ok := urlmap[url]
// 	if !ok {
// 		return "The short url not found"
// 	}
// 	return long
// }

// func geturlHandle(rw http.ResponseWriter, r *http.Request) {
// 	long := geturlFunc(chi.URLParam(r, "id"))
// 	if long == "The short url not found" {
// 		rw.WriteHeader(http.StatusNotFound)
// 		return
// 	}
// 	http.Redirect(rw, r, long, http.StatusTemporaryRedirect)
// 	// rw.WriteHeader(http.StatusTemporaryRedirect)
// 	io.WriteString(rw, long)
// }

// type URL struct {
// 	URL string `json:"url"`
// }

// func makeshortHandle(rw http.ResponseWriter, r *http.Request) {
// 	if r.Host != "localhost:8080" {
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	var url URL
// 	err := json.NewDecoder(r.Body).Decode(&url)
// 	if err != nil {
// 		rw.WriteHeader(http.StatusBadRequest)
// 		rw.Write([]byte(err.Error()))
// 		return
// 	}

// 	rw.WriteHeader(http.StatusCreated)
// 	rw.Write([]byte("http://localhost:8080/" + makeshortFunc(url.URL)))
// }

// // curl -X POST 'http://localhost:8080/' -H "text/plain" -d '{"URL": "abc"}'

// func main() {
// 	// urlmap := make(map[string]string)
// 	r := chi.NewRouter()
// 	r.Post("/", makeshortHandle)
// 	r.Get("/{id}", geturlHandle)
// 	log.Fatal(http.ListenAndServe(":8080", r))
// }

// // k0CA1T3m
