package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/OmarLP/api_module/pkg/authorization"
	"github.com/OmarLP/api_module/pkg/utils"
)

type contextKey string

const UserKey contextKey = "user_claims"

// intercepta peticiones http a rutas protegidas
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		// verificamos que no venga vacío
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "authorization header is required")
			return
		}

		// el formato esperado debe ser "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower((parts[0])) != "bearer" {
			utils.RespondError(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		// validar el token y comprobar expiración
		claims, err := authorization.ValidateToken(parts[1])
		if err != nil {
			// si el token pasó el tiempo establecido, devulve el mensaje de error
			utils.RespondError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// pasar los claims al contexto de la aplicación
		ctx := context.WithValue(r.Context(), UserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
