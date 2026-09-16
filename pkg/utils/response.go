package utils

import (
	"encoding/json"
	"net/http"
)

func RespondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": status,
		"error":  message,
	})
}
