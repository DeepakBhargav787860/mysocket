package securemiddleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("dsfghjklljghjhhjkljfgsgdfhgjhhjhkjlxfcgvhjertertgyhujrtgh")

func GenerateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // token expires in 24 hours
		"iat":     time.Now().Unix(),                     // issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	return token.SignedString(jwtSecret)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("Authorization")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("cookie", cookie)
		// 1. Get the token from the Authorization header
		// authHeader := r.Header.Get("Authorization")
		// log.Println("header", authHeader)
		// if authHeader == "" {
		// 	http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		// 	return
		// }
		tokenString := cookie.Value

		// 2. Extract token from Bearer
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// if tokenString == authHeader {
		// 	http.Error(w, "Invalid token format", http.StatusUnauthorized)
		// 	return
		// }

		// 3. Parse and verify token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure correct signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 4. Token is valid, optionally get user info from claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := claims["user_id"]
			fmt.Println("Authenticated user ID:", userID)
			// You can attach userID to context here
		}

		// 5. Proceed to next handler
		next(w, r)
	}
}
