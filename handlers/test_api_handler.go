package handler

import (
	"encoding/json"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status  string
		Message string
	}

	resp := &response{
		Status:  "ok",
		Message: "Hello World!",
	}
	json.NewEncoder(w).Encode(&resp)
	return
}
