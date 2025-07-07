package securemiddleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Secret key
var jwtSecret = []byte("dsfghjklljghjhhjkljfgsgdfhgjhhjhkjlxfcgvhjertertgyhujrtgh")

// Custom key type for context
type contextKey string

const userIDKey = contextKey("user_id")

// ✅ Generate JWT token
func GenerateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // expires in 24 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ✅ Middleware to verify token and set user ID in context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			log.Println("error reading cookie:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("cookie", cookie)
		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("invalid token:", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// ✅ Get user ID from token
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user UUID format", http.StatusUnauthorized)
			return
		}
		log.Println("uuuuuuuuuuuid", userUUID)
		// ✅ Add userID to context
		ctx := context.WithValue(r.Context(), userIDKey, userUUID)

		// Pass modified request to next handler
		next(w, r.WithContext(ctx))
	}
}

// ✅ Function to extract user ID from context in handler
func GetUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	userID := r.Context().Value(userIDKey)
	if userID == nil {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user ID type")
	}
	return userUUID, nil
}
