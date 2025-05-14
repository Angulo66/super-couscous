package handler

import (
	http_response_encoder "go-basics/pkg/http"
	"net/http"
)

// APIHandler handles requests to the /api endpoint
type APIHandler struct {
	encoder http_response_encoder.ResponseEncoder
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(encoder http_response_encoder.ResponseEncoder) *APIHandler {
	return &APIHandler{encoder: encoder}
}

// ServeHTTP handle HTTP requests to the /api endpoint
func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.encoder.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	// Return API information and available endpoints
	apiInfo := map[string]interface{}{
		"name":        "Go API",
		"version":     "1.0.0",
		"description": "Go API",
		"endpoints": []map[string]string{
			{
				"path":        "/api",
				"description": "API Information",
				"method":      "GET",
			},
		},
	}

	h.encoder.SendSuccess(w, apiInfo, "API Information")
}
