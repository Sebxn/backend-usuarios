package middleware

import (
	"context"
	"net/http"
	"strings"
	"os"
	"github.com/golang-jwt/jwt/v5"
)

// Definir una llave tipo para evitar colisiones con otros contextos
type contextKey struct {
    name string
}

// Llave del contexto para nuestra información del usuario
var UserContextKey = &contextKey{"user"}

func TokenAuthMiddleware(next http.Handler) http.Handler {
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtSecret := os.Getenv("JWT_SECRET_KEY")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Falta token de autorización", http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
			return
		}

		tokenString := splitToken[1]

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil // Reemplazar con tu clave secreta
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Agregar información del token al contexto si es necesario
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
	})
}