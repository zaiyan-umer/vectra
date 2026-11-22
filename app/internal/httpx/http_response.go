package httpx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorPayload struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorPayload `json:"error"`
}

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

type Code string

const (
	ErrCodeNotFound         Code = "not_found"
	ErrCodeInternalError    Code = "internal_error"
	ErrCodeInvalidJSON      Code = "invalid_json"
	ErrCodeValidationFailed Code = "validation_failed"
	ErrCodeUnauthenticated  Code = "unauthenticated"
	ErrCodeForbidden        Code = "forbidden"
	ErrCodeUnavailable      Code = "service_unavailable"
	ErrCodeBadGateway       Code = "bad_gateway"
)

func WriteError(w http.ResponseWriter, httpStatus int, message string, code Code) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorPayload{
			Code:    code,
			Message: message,
		},
	}); err != nil {
		slog.Error("failed to encode error response", "err", err)
	}
}

func WriteJson[T any](w http.ResponseWriter, status int, data T) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(SuccessResponse[T]{Data: data}); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(buf.Bytes())
	return err
}
