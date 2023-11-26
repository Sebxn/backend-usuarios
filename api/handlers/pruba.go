package handlers

import (
    "backend/api/middleware"
    "fmt"
    "github.com/golang-jwt/jwt/v5"
    "net/http"
)

func Prueba(w http.ResponseWriter, r *http.Request) {
    // Recuperar la información del usuario desde el contexto
    claims, ok := r.Context().Value(middleware.UserContextKey).(*jwt.MapClaims)
    if !ok {
        http.Error(w, "Información del usuario no disponible", http.StatusInternalServerError)
        return
    }

    // Ahora tienes tus claims y puedes utilizar la información como necesites
    uid := (*claims)["uid"].(string) // Asegúrate de que este campo existe en tus claims
    fmt.Fprintf(w, "El UID del usuario es: %s", uid)
    
   
}