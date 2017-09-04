package lib

import (
	"encoding/json"
	"net/http"
)

// RenderJSON write data as a json
func RenderJSON(w http.ResponseWriter, data interface{}, err error) {
	if IsInvalid(w, err) {
		return
	}
	response, err := json.Marshal(data)
	if IsInvalid(w, err) {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(response)
}

// IsInvalid returns error status if the err it not nil
func IsInvalid(w http.ResponseWriter, err error) (invalid bool) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return
}
