package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

var logfile *os.File

func InitLogger() {
	var err error
	logfile, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}
	log.SetOutput(logfile)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}
		method := r.Method
		url := r.URL.Path

		log.Printf("[%s] %s - %s %s\n", timestamp, clientIP, method, url)
		next.ServeHTTP(w, r)
	})
}
