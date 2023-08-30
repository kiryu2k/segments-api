package handlers

import (
	"encoding/json"
	"io"
)

type responseError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func WriteJSONError(w io.Writer, status int, msg string) {
	_ = json.NewEncoder(w).Encode(
		&responseError{
			StatusCode: status,
			Message:    msg,
		},
	)
}

func WriteServerError(w io.Writer, status int) {
	WriteJSONError(w, status, "server error")
}
