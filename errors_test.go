package bravesearch

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAPIError tests the APIError type
func TestAPIError(t *testing.T) {
	// Test with wrapped error
	apiErr := &APIError{
		StatusCode: 429,
		Message:    "Too Many Requests",
		Err:        ErrRateLimit,
	}
	
	// Test Error method
	errStr := apiErr.Error()
	assert.Contains(t, errStr, "Too Many Requests")
	assert.Contains(t, errStr, "429")
	assert.Contains(t, errStr, ErrRateLimit.Error())
	
	// Test Unwrap method
	unwrapped := apiErr.Unwrap()
	assert.Equal(t, ErrRateLimit, unwrapped)
	
	// Test without wrapped error
	apiErr = &APIError{
		StatusCode: 404,
		Message:    "Not Found",
	}
	errStr = apiErr.Error()
	assert.Contains(t, errStr, "Not Found")
	assert.Contains(t, errStr, "404")
	assert.NotContains(t, errStr, "nil")
	
	// Test nil unwrap
	unwrapped = apiErr.Unwrap()
	assert.Nil(t, unwrapped)
}

// TestNewAPIError tests creating a new APIError
func TestNewAPIError(t *testing.T) {
	apiErr := NewAPIError(401, "Unauthorized", ErrUnauthorized)
	assert.Equal(t, 401, apiErr.StatusCode)
	assert.Equal(t, "Unauthorized", apiErr.Message)
	assert.Equal(t, ErrUnauthorized, apiErr.Err)
}

// TestNewHTTPError tests creating an APIError from HTTP response
func TestNewHTTPError(t *testing.T) {
	// Test 401 Unauthorized
	resp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Status:     "401 Unauthorized",
	}
	apiErr := NewHTTPError(resp)
	assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	assert.Equal(t, "401 Unauthorized", apiErr.Message)
	assert.Equal(t, ErrUnauthorized, apiErr.Err)
	
	// Test 403 Forbidden
	resp = &http.Response{
		StatusCode: http.StatusForbidden,
		Status:     "403 Forbidden",
	}
	apiErr = NewHTTPError(resp)
	assert.Equal(t, http.StatusForbidden, apiErr.StatusCode)
	assert.Equal(t, ErrForbidden, apiErr.Err)
	
	// Test 404 Not Found
	resp = &http.Response{
		StatusCode: http.StatusNotFound,
		Status:     "404 Not Found",
	}
	apiErr = NewHTTPError(resp)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Equal(t, ErrNotFound, apiErr.Err)
	
	// Test 429 Too Many Requests
	resp = &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Status:     "429 Too Many Requests",
	}
	apiErr = NewHTTPError(resp)
	assert.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
	assert.Equal(t, ErrRateLimit, apiErr.Err)
	
	// Test 500 Server Error
	resp = &http.Response{
		StatusCode: http.StatusInternalServerError,
		Status:     "500 Internal Server Error",
	}
	apiErr = NewHTTPError(resp)
	assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
	assert.Equal(t, ErrServerError, apiErr.Err)
	
	// Test unknown error code
	resp = &http.Response{
		StatusCode: 418,
		Status:     "418 I'm a teapot",
	}
	apiErr = NewHTTPError(resp)
	assert.Equal(t, 418, apiErr.StatusCode)
	assert.Contains(t, apiErr.Err.Error(), "unexpected status code: 418")
}

// TestIsRateLimitError tests the rate limit error detection
func TestIsRateLimitError(t *testing.T) {
	// Test direct error
	assert.True(t, IsRateLimitError(ErrRateLimit))
	
	// Test wrapped in APIError
	apiErr := &APIError{
		StatusCode: 429,
		Message:    "Rate limit exceeded",
		Err:        ErrRateLimit,
	}
	assert.True(t, IsRateLimitError(apiErr))
	
	// Test with status code but different error
	apiErr = &APIError{
		StatusCode: 429,
		Message:    "Rate limit exceeded",
		Err:        errors.New("custom error"),
	}
	assert.True(t, IsRateLimitError(apiErr))
	
	// Test with different status code
	apiErr = &APIError{
		StatusCode: 401,
		Message:    "Unauthorized",
		Err:        ErrUnauthorized,
	}
	assert.False(t, IsRateLimitError(apiErr))
	
	// Test with non-rate limit error
	assert.False(t, IsRateLimitError(ErrUnauthorized))
	assert.False(t, IsRateLimitError(errors.New("some other error")))
}

// TestIsAuthError tests the authentication error detection
func TestIsAuthError(t *testing.T) {
	// Test direct errors
	assert.True(t, IsAuthError(ErrUnauthorized))
	assert.True(t, IsAuthError(ErrInvalidAPIKey))
	
	// Test wrapped in APIError with matching error
	apiErr := &APIError{
		StatusCode: 401,
		Message:    "Unauthorized",
		Err:        ErrUnauthorized,
	}
	assert.True(t, IsAuthError(apiErr))
	
	// Test with status code but different error
	apiErr = &APIError{
		StatusCode: 401,
		Message:    "Unauthorized",
		Err:        errors.New("custom error"),
	}
	assert.True(t, IsAuthError(apiErr))
	
	// Test with different status code but matching error
	apiErr = &APIError{
		StatusCode: 400,
		Message:    "Bad Request",
		Err:        ErrUnauthorized,
	}
	assert.True(t, IsAuthError(apiErr))
	
	// Test with non-auth error
	assert.False(t, IsAuthError(ErrRateLimit))
	assert.False(t, IsAuthError(errors.New("some other error")))
	
	// Test with non-auth error and non-auth status code
	apiErr = &APIError{
		StatusCode: 429,
		Message:    "Rate limit exceeded",
		Err:        ErrRateLimit,
	}
	assert.False(t, IsAuthError(apiErr))
}

// TestIsServerError tests the server error detection
func TestIsServerError(t *testing.T) {
	// Test direct error
	assert.True(t, IsServerError(ErrServerError))
	
	// Test wrapped in APIError with matching error
	apiErr := &APIError{
		StatusCode: 500,
		Message:    "Internal Server Error",
		Err:        ErrServerError,
	}
	assert.True(t, IsServerError(apiErr))
	
	// Test with status code but different error
	apiErr = &APIError{
		StatusCode: 502,
		Message:    "Bad Gateway",
		Err:        errors.New("custom error"),
	}
	assert.True(t, IsServerError(apiErr))
	
	// Test with different status code but matching error
	apiErr = &APIError{
		StatusCode: 400,
		Message:    "Bad Request",
		Err:        ErrServerError,
	}
	assert.True(t, IsServerError(apiErr))
	
	// Test with non-server error
	assert.False(t, IsServerError(ErrRateLimit))
	assert.False(t, IsServerError(errors.New("some other error")))
	
	// Test with non-server error and non-server status code
	apiErr = &APIError{
		StatusCode: 429,
		Message:    "Rate limit exceeded",
		Err:        ErrRateLimit,
	}
	assert.False(t, IsServerError(apiErr))
}
