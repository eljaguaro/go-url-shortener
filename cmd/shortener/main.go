package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
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
			"uri", r.RequestURI,
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
	r.Post("/", WithLogging((func(rw http.ResponseWriter, r *http.Request) { makeshortHandle(rw, r, *surladdr) })))
	r.Get("/{id}", WithLogging(geturlHandle))

	addr := "127.0.0.1:8080"
	sugar.Infow(
		"Starting server",
		"addr", addr,
	)
	if err := http.ListenAndServe(addr, nil); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
