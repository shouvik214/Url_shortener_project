package handler

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

// func writeSuccess(w http.ResponseWriter, status int, message string) {
// 	writeJSON(w, status, SuccessResponse{
// 		Message: message,
// 	})
// }
