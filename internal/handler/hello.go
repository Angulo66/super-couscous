package handler

import (
	http_response_encoder "go-basics/pkg/http"
	"net/http"
)

type HelloHandler struct {
	encoder http_response_encoder.ResponseEncoder
}

func NewHelloHandler(encoder http_response_encoder.ResponseEncoder) *HelloHandler {
	return &HelloHandler{encoder: encoder}
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.encoder.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	h.encoder.SendSuccess(w, map[string]string{"message": "Hello World!"}, "Hello, World!")
}
