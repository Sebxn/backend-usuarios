package middleware

import (
    "net/http"
)

// CorsMiddleware maneja las solicitudes CORS con una lista de orígenes permitidos.
func CorsMiddleware(allowedOrigins []string) func(handler http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")

            if isOriginAllowed(origin, allowedOrigins) {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,PATCH")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Incluye Authorization aquí
            }

            if r.Method == "OPTIONS" {
                if isOriginAllowed(origin, allowedOrigins) {
                    w.WriteHeader(http.StatusOK)
                } else {
                    w.WriteHeader(http.StatusForbidden)
                }
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

func DisableCors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Permitir cualquier origen
        w.Header().Set("Access-Control-Allow-Origin", "*")

        // Permitir cualquier método
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")

        // Permitir cualquier encabezado
        w.Header().Set("Access-Control-Allow-Headers", "*")

        // Permitir credenciales - ten en cuenta que no puedes usar credenciales con un origen de '*'
        // w.Header().Set("Access-Control-Allow-Credentials", "true")

        // Si es una solicitud de pre-vuelo, detener aquí
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Servir la siguiente middleware o handler en la cadena
        next.ServeHTTP(w, r)
    })
}

// isOriginAllowed verifica si el origen de la solicitud está en la lista de orígenes permitidos.
func isOriginAllowed(origin string, allowedOrigins []string) bool {
    for _, allowedOrigin := range allowedOrigins {
        if allowedOrigin == origin {
            return true
        }
    }
    return false
}
