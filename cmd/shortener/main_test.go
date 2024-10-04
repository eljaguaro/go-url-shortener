package main

import (
	"net/http"

	// "net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *resty.Client, method string,
	path string, body string) *resty.Response {
	if method == "POST" {
		resp, err := ts.R().
			SetHeader("Content-Type", "text/plain; charset=UTF-8").
			SetBody(body).
			Post(path)
		require.NoError(t, err)
		return resp
	} else {
		resp, err := ts.R().
			SetHeader("Content-Type", "text/plain; charset=UTF-8").
			SetBody(body).
			Get(path)
		require.NoError(t, err)
		return resp
	}
}

func TestMakeshort(t *testing.T) {
	ts := resty.New().SetBaseURL("http://localhost:8080")
	type etal struct {
		method string
		url    string
		body   string
		status int
		geturl string
	}
	var testTable = []etal{
		// {"POST", "http://localhost:8080/", "", http.StatusCreated, ""},
		{"POST", "http://localhost:8080/", "https://practicum.ru/", http.StatusCreated, ""},
		{"POST", "http://localhost:8080/", "https://yandex.ru/", http.StatusCreated, ""},
	}
	var shorts []string
	for _, v := range testTable {
		resp := testRequest(t, ts, v.method, v.url, v.body)
		assert.Equal(t, v.status, resp.StatusCode())
		shorts = append(shorts, string(resp.Body()))
	}

	testTable = append(testTable, etal{"GET", "http://localhost:8080/" + shorts[0], "", http.StatusTemporaryRedirect, "https://practicum.ru/"})
	testTable = append(testTable, etal{"GET", "http://localhost:8080/" + shorts[1], "", http.StatusTemporaryRedirect, "https://yandex.ru/"})
	testTable = append(testTable, etal{"GET", "http://localhost:8080/" + "A" + shorts[1], "", http.StatusNotFound, ""})
	for _, v := range testTable[2:] {
		resp := testRequest(t, ts, v.method, v.url, v.body)
		assert.Equal(t, v.status, resp.StatusCode())
		assert.Equal(t, v.geturl, string(resp.Body()))
	}
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
// 	if r.Host != "localhost:8080" {
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	long := geturlFunc(chi.URLParam(r, "id"))
// 	if long == "The short url not found" {
// 		rw.WriteHeader(http.StatusNotFound)
// 		return
// 	}
// 	rw.WriteHeader(http.StatusTemporaryRedirect)
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
// 	rw.Write([]byte(makeshortFunc(url.URL)))
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
