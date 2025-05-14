package http_response_encoder

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ResponseEncoder defines the interface for encoding HTTP responses allows for extension
type ResponseEncoder interface {
	SendSuccess(w http.ResponseWriter, data interface{}, message string)
	SendError(w http.ResponseWriter, statusCode int, errorMsg string)
}

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSONEncoder implements ResponseEncoder for JSON responses
type JSONEncoder struct{}

// NewJSONEncoder creates a new JSON response encoder
func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

// sendSuccess sends a success JSON response
func (e *JSONEncoder) SendSuccess(w http.ResponseWriter, data interface{}, message string) {
	sendJson(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// sendError sends an error JSON response
func (e *JSONEncoder) SendError(w http.ResponseWriter, statusCode int, errorMsg string) {
	sendJson(w, statusCode, Response{Success: false, Error: errorMsg})
}

// sendJson is a helper function to send JSON responses
func sendJson(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println(err)
	}
}
