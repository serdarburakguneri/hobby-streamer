package shared

import (
	"encoding/json"
	"net/http"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set(constants.HeaderContentType, "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, statusCode int, error string) {
	w.Header().Set(constants.HeaderContentType, "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Error: error,
	}

	json.NewEncoder(w).Encode(response)
}
