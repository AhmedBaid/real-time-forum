package config

import (
	"encoding/json"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, status_code int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	json.NewEncoder(w).Encode(message)
}
