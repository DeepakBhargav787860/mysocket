package cros

import (
	"log"
	"net/http"
)

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Println("origin", origin)
		// if origin == "https://chat-sj5k-deepaks-projects-a5241927.vercel.app" || origin == "http://localhost:5173" {

		// }
		w.Header().Set("Access-Control-Allow-Origin", "https://chat-sj5k-deepaks-projects-a5241927.vercel.app")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
