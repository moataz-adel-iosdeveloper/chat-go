package middlewares

import (
	"chat-go/internal/config"
	userRepository "chat-go/internal/repositories/user"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeards := r.Header.Get("Authorization")
		if authHeards == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Here you can add more logic to validate the token
		parts := strings.Split(authHeards, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(parts[1])
		// log.Printf("Auth token received: %s", token)

		// Validate JWT signature and extract user_id claim
		userID, err := config.ValidateToken(token)
		if err != nil {
			log.Printf("Invalid JWT token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Ensure the token matches what's stored for the user
		user, err := userRepository.FindUserByID(userID)
		if err != nil || user == nil {
			log.Printf("User not found for id: %s, err: %v", userID, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if strings.TrimSpace(user.Token) != token {
			log.Printf("Token mismatch for user %s", userID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func AuthMiddlewareSocket(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Validate JWT signature and extract user_id claim
		userID, err := config.ValidateToken(token)
		if err != nil {
			log.Printf("Invalid JWT token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Ensure the token matches what's stored for the user
		user, err := userRepository.FindUserByID(userID)
		if err != nil || user == nil {
			log.Printf("User not found for id: %s, err: %v", userID, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if strings.TrimSpace(user.Token) != token {
			log.Printf("Token mismatch for user %s", userID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func WebAuthMiddleware(store config.Store) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, _ := store.Get(r, "auth-session")

			if sess.Values["user_id"] == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
