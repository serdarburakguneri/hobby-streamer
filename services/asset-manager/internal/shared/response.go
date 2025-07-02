package shared

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data interface{}, errMsg string) {
	resp := map[string]interface{}{
		"data":  data,
		"error": nil,
	}
	if errMsg != "" {
		resp["data"] = nil
		resp["error"] = errMsg
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}