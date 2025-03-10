package bravesearch

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClient tests the client creation
func TestNewClient(t *testing.T) {
	// Test with missing API key
	client, err := NewClient("")
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrMissingAPIKey, err)

	// Test with valid API key
	client, err = NewClient("test-api-key")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "test-api-key", client.config.APIKey)
	assert.Equal(t, BaseURL, client.config.BaseURL)
	assert.Equal(t, DefaultUserAgent, client.config.UserAgent)
	assert.Equal(t, time.Duration(DefaultTimeout)*time.Second, client.config.Timeout)

	// Test with options
	client, err = NewClient("test-api-key",
		WithTimeout(60),
		WithUserAgent("custom-agent"),
		WithDefaultCountry("UK"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "test-api-key", client.config.APIKey)
	assert.Equal(t, 60*time.Second, client.config.Timeout)
	assert.Equal(t, "custom-agent", client.config.UserAgent)
	assert.Equal(t, "UK", client.config.DefaultCountry)
}

// TestWebSearch tests the web search functionality
func TestWebSearch(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/res/v1/web/search", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get(HeaderSubscriptionToken))
		assert.Equal(t, MIMETypeJSON, r.Header.Get(HeaderAccept))
		assert.Equal(t, DefaultUserAgent, r.Header.Get(HeaderUserAgent))

		// Check query parameters
		q := r.URL.Query().Get("q")
		assert.Equal(t, "go programming", q)

		// Set rate limit headers
		w.Header().Set(HeaderRateLimitLimit, "10, 15000")
		w.Header().Set(HeaderRateLimitRemaining, "9, 14999")
		w.Header().Set(HeaderRateLimitReset, "1, 1419704")

		// Return test response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Read test data
		data, err := os.ReadFile("testdata/web_search_response.json")
		require.NoError(t, err)
		_, err = w.Write(data)
		require.NoError(t, err)
	}))
	defer server.Close()

	// Create client with mock server
	client, err := NewClient("test-api-key", WithBaseURL(server.URL+"/res/v1"))
	require.NoError(t, err)
	require.NotNil(t, client)

	// Perform search
	ctx := context.Background()
	response, err := client.WebSearch(ctx, "go programming", nil)
	require.NoError(t, err)
	require.NotNil(t, response)

	// Verify response
	assert.Equal(t, "search", response.Type)
	assert.NotNil(t, response.Web)
	assert.Equal(t, 3, len(response.Web.Results))
	assert.Equal(t, "The Go Programming Language", response.Web.Results[0].Title)
	assert.Equal(t, "https://go.dev/", response.Web.Results[0].URL)
	assert.Contains(t, response.Web.Results[0].Description, "Go is an open source programming language")
}

// TestWebSearchEmptyQuery tests the validation for empty queries
func TestWebSearchEmptyQuery(t *testing.T) {
	client, err := NewClient("test-api-key")
	require.NoError(t, err)
	require.NotNil(t, client)

	// Test with empty query
	response, err := client.WebSearch(context.Background(), "", nil)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrEmptyQuery, err)
}

// TestWebSearchError tests error handling
func TestWebSearchError(t *testing.T) {
	// Setup test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "Invalid API key"}`))
	}))
	defer server.Close()

	// Create client with mock server
	client, err := NewClient("test-api-key", WithBaseURL(server.URL+"/res/v1"))
	require.NoError(t, err)
	require.NotNil(t, client)

	// Perform search
	response, err := client.WebSearch(context.Background(), "go programming", nil)
	assert.Error(t, err)
	assert.Nil(t, response)

	// Check error type
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	assert.Equal(t, ErrUnauthorized, apiErr.Unwrap())

	// Check using helper function
	assert.True(t, IsAuthError(err))
	assert.False(t, IsRateLimitError(err))
}

// TestRateLimitError tests rate limit error detection
func TestRateLimitError(t *testing.T) {
	// Setup test server that returns a rate limit error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error": "Rate limit exceeded"}`))
	}))
	defer server.Close()

	// Create client with mock server
	client, err := NewClient("test-api-key", WithBaseURL(server.URL+"/res/v1"))
	require.NoError(t, err)
	require.NotNil(t, client)

	// Perform search
	response, err := client.WebSearch(context.Background(), "go programming", nil)
	assert.Error(t, err)
	assert.Nil(t, response)

	// Check error type
	assert.True(t, IsRateLimitError(err))
	assert.False(t, IsServerError(err))

	// Check error details
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
}

// TestParseRateLimitHeaders tests parsing of rate limit headers
func TestParseRateLimitHeaders(t *testing.T) {
	// Setup mock response
	resp := &http.Response{
		Header: http.Header{},
	}
	resp.Header.Set(HeaderRateLimitLimit, "10, 15000")
	resp.Header.Set(HeaderRateLimitRemaining, "9, 14999")
	resp.Header.Set(HeaderRateLimitReset, "1, 1419704")

	// Create client
	client, err := NewClient("test-api-key")
	require.NoError(t, err)
	require.NotNil(t, client)

	// Parse headers
	rateLimit := client.parseRateLimitHeaders(resp)
	assert.NotNil(t, rateLimit)
	assert.Equal(t, 10, rateLimit.Limit)
	assert.Equal(t, 9, rateLimit.Remaining)
	assert.Equal(t, 1, rateLimit.Reset)
}

// TestMakeRequestWithRetries tests the retry mechanism
func TestMakeRequestWithRetries(t *testing.T) {
	// Setup test server that fails twice then succeeds
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++

		// First two attempts return 503 error
		if attempts <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error": "Service unavailable"}`))
			return
		}

		// Third attempt succeeds
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"type": "search", "web": {"results": []}}`))
	}))
	defer server.Close()

	// Create client with options for retries and pointing to our mock server
	client, err := NewClient("test-api-key",
		WithBaseURL(server.URL),
		WithRetries(2))
	require.NoError(t, err)
	require.NotNil(t, client)

	// Make a request to the test server
	var response WebSearchResponse
	err = client.makeRequest(context.Background(), http.MethodGet, server.URL, nil, &response)
	assert.NoError(t, err)
	assert.Equal(t, 3, attempts) // Original request + 2 retries = 3 attempts total
}

// TestBuildRequestURL tests URL building with query parameters
func TestBuildRequestURL(t *testing.T) {
	client, err := NewClient("test-api-key")
	require.NoError(t, err)
	require.NotNil(t, client)

	// Test with minimal parameters
	params := &WebSearchParams{
		Query: "test query",
	}
	url, err := client.buildRequestURL(WebSearchEndpoint, params)
	assert.NoError(t, err)
	assert.Contains(t, url, WebSearchEndpoint)
	assert.Contains(t, url, "q=test+query")

	// Test with all parameters
	params = &WebSearchParams{
		Query:           "test query",
		Country:         "JP",
		SearchLang:      "ja",
		UILang:          "ja-JP",
		Count:           10,
		Offset:          2,
		SafeSearch:      SafeSearchStrict,
		Freshness:       FreshnessWeek,
		TextDecorations: true,
		Spellcheck:      false,
		ResultFilter:    ResultFilterNews,
		Goggles:         "custom-goggle",
		Units:           UnitMetric,
		ExtraSnippets:   true,
		Summary:         true,
	}
	url, err = client.buildRequestURL(WebSearchEndpoint, params)
	assert.NoError(t, err)
	assert.Contains(t, url, "q=test+query")
	assert.Contains(t, url, "country=JP")
	assert.Contains(t, url, "search_lang=ja")
	assert.Contains(t, url, "ui_lang=ja-JP")
	assert.Contains(t, url, "count=10")
	assert.Contains(t, url, "offset=2")
	assert.Contains(t, url, "safesearch=strict")
	assert.Contains(t, url, "freshness=pw")
	assert.Contains(t, url, "text_decorations=true")
	assert.Contains(t, url, "spellcheck=false")
	assert.Contains(t, url, "result_filter=news")
	assert.Contains(t, url, "goggles=custom-goggle")
	assert.Contains(t, url, "units=metric")
	assert.Contains(t, url, "extra_snippets=true")
	assert.Contains(t, url, "summary=true")
}

// loadTestData loads test response data
func loadTestData(t *testing.T, path string) *WebSearchResponse {
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var response WebSearchResponse
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	return &response
}
