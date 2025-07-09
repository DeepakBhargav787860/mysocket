package cros

import "net/http"

func EnableCORS(next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{
		"http://127.0.0.1:5173":                      true,
		"http://localhost:5173":                      true,
		"https://socket-application.vercel.app":      true,
		"https://socket-application-mp9v.vercel.app": true,
		"https://chat-sj5k.vercel.app":               true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
