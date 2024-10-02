package main

import (
	"encoding/json"
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
	long := geturlFunc(chi.URLParam(r, "id"))
	if long == "The short url not found" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusTemporaryRedirect)
	io.WriteString(rw, long)
}

type Url struct {
	Url string `json:"url"`
}

func makeshortHandle(rw http.ResponseWriter, r *http.Request) {
	var url Url
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(makeshortFunc(url.Url)))
}

// curl -X POST 'http://localhost:8080/' -H "text/plain" -d '{"URL": "abc"}'

func main() {
	// urlmap := make(map[string]string)
	r := chi.NewRouter()
	r.Post("/", makeshortHandle)
	r.Get("/{id}", geturlHandle)
	log.Fatal(http.ListenAndServe(":8080", r))
}

// k0CA1T3m
