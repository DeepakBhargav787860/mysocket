package response

import (
	"encoding/json"
	"net/http"
)

func MessagePassed(w http.ResponseWriter, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "success",
		"data":    message,
	})
}
