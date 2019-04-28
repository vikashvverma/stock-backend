package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response represent the response to be sent for an API call.
type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Errors  *Error      `json:"errors,omitempty"`
}

// Error represents an error to be sent to the client.
type Error struct {
	Reason string `json:"reason,omitempty"`
}

// Send writes a successful response to the given http.ResponseWriter.
func (s Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		return fmt.Errorf("send: unable to write JSON response: %s", err)
	}

	return nil
}

// ServerError writes a server error response to the given http.ResponseWriter.
func (s Response) ServerError(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		return fmt.Errorf("serverError: could not write JSON response: %s", err)
	}

	return nil
}

// ClientError writes a client error response to the given http.ResponseWriter.
func (s Response) ClientError(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		return fmt.Errorf("clientError: could not write JSON response: %s", err)
	}

	return nil
}



func (e Error) Forbidden(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)

	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		return fmt.Errorf("unauthorized: could not write JSON response: %s", err)
	}

	return nil
}
