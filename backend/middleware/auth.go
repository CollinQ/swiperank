package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

// Define custom context key types to prevent collisions
type contextKey string

const (
	userIDKey contextKey = "user_id"
	userKey   contextKey = "user"
)

// Initialize Clerk globally with an API key
func InitClerk(apiKey string) {
	clerk.SetKey(apiKey) // Set Clerk API key globally
	log.Println("✅ Clerk API Key initialized")
}

// Middleware to handle authentication using Clerk
func AuthMiddleware(next http.Handler) http.Handler {
	return clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("🟢 AuthMiddleware: Middleware is executing!")

		// Log Authorization header
		authHeader := r.Header.Get("Authorization")
		log.Println("🔹 Received Authorization Header:", authHeader)

		ctx := r.Context()

		// Extract Clerk session claims
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			log.Println("❌ AuthMiddleware: No valid Clerk session found")
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		log.Println("✅ AuthMiddleware: Clerk session found for user:", claims.Subject)

		// Fetch user details from Clerk using their ID (claims.Subject)
		usr, err := user.Get(ctx, claims.Subject)
		if err != nil {
			log.Println("❌ AuthMiddleware: Error fetching user from Clerk:", err)
			http.Error(w, `{"error": "Failed to retrieve user"}`, http.StatusInternalServerError)
			return
		}
		if usr == nil {
			log.Println("❌ AuthMiddleware: User does not exist")
			http.Error(w, `{"error": "User does not exist"}`, http.StatusNotFound)
			return
		}

		log.Println("✅ AuthMiddleware: Authenticated User ID:", claims.Subject)

		// Attach the user ID and user object to the request context
		ctx = context.WithValue(ctx, userIDKey, claims.Subject)
		ctx = context.WithValue(ctx, userKey, usr)
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}
