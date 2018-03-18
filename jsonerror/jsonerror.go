package jsonerror

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

// JsonError works like http.Error but uses our response
// struct as the body of the response. Like http.Error
// you will still need to call a naked return in the http handler
func Error(w http.ResponseWriter, r *Response, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	b, _ := json.Marshal(r)

	w.Write(b)
}
