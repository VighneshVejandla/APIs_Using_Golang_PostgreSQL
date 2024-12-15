package handlers

import "net/http"

// ErrorHandler handles error responses
func ErrorHandler(w http.ResponseWriter, message string, code int) {
    http.Error(w, message, code)
}
