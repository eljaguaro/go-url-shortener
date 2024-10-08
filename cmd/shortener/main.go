package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/caarlos0/env/v6"
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

func makeshortHandle(rw http.ResponseWriter, r *http.Request, surladdr string) {
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
	rw.Write([]byte(surladdr + "/" + makeshortFunc(url.URL)))
}

// func rhandle(rw http.ResponseWriter, r *http.Request) {
// 	http.Redirect(rw, r, "http://abc.ru/", 309)

// }

// curl -X POST 'http://localhost:8080/' -H "text/plain" -d '{"URL": "abc"}'
type Config struct {
	SERVER_ADDRESS *string `env:"SERVER_ADDRESS"`
	BASE_URL       *string `env:"BASE_URL"`
}

func main() {
	// urlmap := make(map[string]string)
	// r.Get("/fw", rhandle)
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	run := flag.String("a", "localhost:8080", "адрес запуска http-сервера")
	surladdr := flag.String("b", "http://localhost:8080", "базовый адрес результирующего URL")
	flag.Parse()
	fmt.Println("address to run the server:", *run)
	fmt.Println("server address and shorturl", *surladdr)
	if *cfg.SERVER_ADDRESS != "" {
		run = cfg.SERVER_ADDRESS
	}
	if *cfg.BASE_URL != "" {
		surladdr = cfg.BASE_URL
	}
	port := strings.Split(*run, ":")[1]
	r := chi.NewRouter()
	r.Post("/", func(rw http.ResponseWriter, r *http.Request) { makeshortHandle(rw, r, *surladdr) })
	r.Get("/{id}", geturlHandle)
	log.Fatal(http.ListenAndServe(":"+port, r))
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
