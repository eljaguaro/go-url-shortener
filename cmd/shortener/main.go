package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/thanhpk/randstr"
	"go.uber.org/zap"
)

var urlmap = make(map[string]string)

func makeshortFunc(long string) string {
	short := randstr.String(8)
	urlmap[short] = long
	return short
}

func geturlFunc(url string) string {
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
	http.Redirect(rw, r, long, http.StatusTemporaryRedirect)
}

type URL struct {
	URL string
}

func makeshortHandle(rw http.ResponseWriter, r *http.Request, surladdr string) {
	var url URL
	b, err := io.ReadAll(r.Body)
	url.URL = string(b)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "text/plain")
	rw.Write([]byte(surladdr + "/" + makeshortFunc(url.URL)))
}

func makeshortjsonHandle(rw http.ResponseWriter, r *http.Request, surladdr string) {
	var url URL
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &url); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["result"] = surladdr + "/" + makeshortFunc(url.URL)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Write(jsonResp)
}

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		sugar.Infoln(
			"method", r.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}
	return http.HandlerFunc(logFn)
}

type Config struct {
	serverAddress *string `env:"SERVER_ADDRESS"`
	baseURL       *string `env:"BASE_URL"`
}

var sugar zap.SugaredLogger

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}
type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

var acceptedHeaders = []string{"application/json", "text/html"}

// func gzipHandle(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || !slices.Contains(acceptedHeaders, r.Header.Get("Content-Type")) {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		// создаём gzip.Writer поверх текущего w
// 		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
// 		if err != nil {
// 			io.WriteString(w, err.Error())
// 			return
// 		}
// 		defer gz.Close()

//			w.Header().Set("Content-Encoding", "gzip")
//			// передаём обработчику страницы переменную типа gzipWriter для вывода данных
//			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
//		})
//	}
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w
		if slices.Contains(acceptedHeaders, r.Header.Get("Content-Type")) {

			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				cw := newCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close()
			}
		}

		h.ServeHTTP(ow, r)
	}
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	run := flag.String("a", "localhost:8080", "адрес запуска http-сервера")
	surladdr := flag.String("b", "http://localhost:8080", "базовый адрес результирующего URL")
	flag.Parse()
	fmt.Println("address to run the server:", run)
	fmt.Println("server address and shorturl", surladdr)
	if cfg.serverAddress != nil {
		run = cfg.serverAddress
	}
	if cfg.baseURL != nil {
		surladdr = cfg.baseURL
	}
	port := strings.Split(*run, ":")[1]
	r := chi.NewRouter()
	r.Post("/", WithLogging(gzipMiddleware((func(rw http.ResponseWriter, r *http.Request) { makeshortHandle(rw, r, *surladdr) }))))
	r.Get("/{id}", WithLogging(gzipMiddleware(geturlHandle)))
	r.Post("/api/shorten", WithLogging(gzipMiddleware((func(rw http.ResponseWriter, r *http.Request) { makeshortjsonHandle(rw, r, *surladdr) }))))
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// curl --header "Content-Type: application/json" \
//    --request POST \
//    --data '{"username":"xyz","password":"xyz"}' \
//    http://localhost:8080/api/shorten
