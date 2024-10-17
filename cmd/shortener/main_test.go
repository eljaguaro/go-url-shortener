package main

import (
	"io"
	"net/http"
	"strings"

	// "net/http/httptest"
	"fmt"
	"testing"

	// "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *http.Client, method string,
	path string, body string) *http.Response {
	if method == "POST" {
		req, err := http.NewRequest(method, path, strings.NewReader(""))
		require.NoError(t, err)
		req.Header.Add("Content-type", `"text/plain"`)
		resp, err := ts.Do(req)
		require.NoError(t, err)
		return resp
	} else {
		req, err := http.NewRequest(method, path, nil)
		require.NoError(t, err)
		req.Header.Add("Content-type", `"text/plain"`)
		resp, err := ts.Do(req)
		require.NoError(t, err)
		return resp
	}
}

var shorts []string

func TestShortener(t *testing.T) {
	// quit := make(chan bool)
	go main()
	ts := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	type etalpost struct {
		method string

		url    string
		body   string
		status int
		geturl string
	}
	var testTable = []etalpost{
		{"POST", "http://localhost:8080/", "http://abc.com", http.StatusCreated, ""},
		{"POST", "http://localhost:8080/", "http://cba.com", http.StatusCreated, ""},
	}
	// var shorts []string
	for _, v := range testTable {
		resp := testRequest(t, ts, v.method, v.url, v.body)

		assert.Equal(t, v.status, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		shorts = append(shorts, string(b))
		resp.Body.Close()
	}
	type etalget struct {
		method string
		url    string
		status int
		scheme string
		host   string
	}

	fmt.Println(shorts)
	var testTable2 = []etalget{}
	testTable2 = append(testTable2, etalget{"GET", shorts[0], 307, "http", "abc.com"})
	testTable2 = append(testTable2, etalget{"GET", shorts[1], 307, "http", "cba.com"})
	fmt.Println(testTable)
	for _, v := range testTable2 {
		fmt.Println(testTable)
		resp := testRequest(t, ts, v.method, v.url, "")
		assert.Equal(t, v.status, resp.StatusCode)
		resp.Body.Close()
	}
}
