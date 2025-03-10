package bravesearch

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrMissingAPIKey is returned when the API key is missing
	ErrMissingAPIKey = errors.New("missing API key")

	// ErrInvalidAPIKey is returned when the API key is invalid
	ErrInvalidAPIKey = errors.New("invalid API key")

	// ErrRateLimit is returned when the API rate limit is exceeded
	ErrRateLimit = errors.New("rate limit exceeded")

	// ErrUnauthorized is returned when the API returns a 401 Unauthorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when the API returns a 403 Forbidden
	ErrForbidden = errors.New("forbidden")

	// ErrNotFound is returned when the API returns a 404 Not Found
	ErrNotFound = errors.New("not found")

	// ErrServerError is returned when the API returns a 5xx error
	ErrServerError = errors.New("server error")

	// ErrInvalidResponse is returned when the API returns an invalid response
	ErrInvalidResponse = errors.New("invalid response")

	// ErrInvalidParameters is returned when invalid parameters are provided
	ErrInvalidParameters = errors.New("invalid parameters")

	// ErrQueryTooLong is returned when the query is too long (>400 chars or >50 words)
	ErrQueryTooLong = errors.New("query too long (max 400 chars or 50 words)")

	// ErrEmptyQuery is returned when an empty query is provided
	ErrEmptyQuery = errors.New("query cannot be empty")
)

// APIError represents an error returned by the Brave Search API
type APIError struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
	Err        error  `json:"error,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("brave search API error: %s (status: %d): %s", e.Message, e.StatusCode, e.Err.Error())
	}
	return fmt.Sprintf("brave search API error: %s (status: %d)", e.Message, e.StatusCode)
}

// Unwrap returns the wrapped error
func (e *APIError) Unwrap() error {
	return e.Err
}

// NewAPIError creates a new APIError
func NewAPIError(statusCode int, message string, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

// NewHTTPError creates a new APIError from an HTTP response
func NewHTTPError(resp *http.Response) *APIError {
	var err error
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		err = ErrUnauthorized
	case http.StatusForbidden:
		err = ErrForbidden
	case http.StatusNotFound:
		err = ErrNotFound
	case http.StatusTooManyRequests:
		err = ErrRateLimit
	default:
		if resp.StatusCode >= 500 {
			err = ErrServerError
		} else {
			err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    resp.Status,
		Err:        err,
	}
}

// IsRateLimitError checks if the error is a rate limit error
func IsRateLimitError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return errors.Is(apiErr.Err, ErrRateLimit) || apiErr.StatusCode == http.StatusTooManyRequests
	}
	return errors.Is(err, ErrRateLimit)
}

// IsAuthError checks if the error is an authentication error
func IsAuthError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return errors.Is(apiErr.Err, ErrUnauthorized) || 
			   errors.Is(apiErr.Err, ErrInvalidAPIKey) ||
			   apiErr.StatusCode == http.StatusUnauthorized
	}
	return errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrInvalidAPIKey)
}

// IsServerError checks if the error is a server error
func IsServerError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return errors.Is(apiErr.Err, ErrServerError) || apiErr.StatusCode >= 500
	}
	return errors.Is(err, ErrServerError)
}
