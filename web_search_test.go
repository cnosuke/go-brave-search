package bravesearch

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMockServer sets up a mock server for testing
func setupMockServer(t *testing.T) (*httptest.Server, *Client) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request is valid
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/res/v1/web/search", r.URL.Path)

		// Extract query parameters
		query := r.URL.Query()
		q := query.Get("q")
		require.NotEmpty(t, q)

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Load and return test data
		data, err := os.ReadFile("testdata/web_search_response.json")
		require.NoError(t, err)
		_, err = w.Write(data)
		require.NoError(t, err)
	}))

	// Create client with mock server
	client, err := NewClient("test-api-key", WithBaseURL(server.URL+"/res/v1"))
	require.NoError(t, err)

	return server, client
}

// TestNewWebSearchParams tests the creation of default search parameters
func TestNewWebSearchParams(t *testing.T) {
	params := NewWebSearchParams()
	assert.NotNil(t, params)
	assert.Equal(t, DefaultCount, params.Count)
	assert.Equal(t, DefaultOffset, params.Offset)
	assert.Equal(t, DefaultSafeSearch, params.SafeSearch)
	assert.Equal(t, DefaultTextDecor, params.TextDecorations)
	assert.Equal(t, DefaultSpellCheck, params.Spellcheck)
}

// TestWebSearchWithCountry tests the search with country helper function
func TestWebSearchWithCountry(t *testing.T) {
	server, client := setupMockServer(t)
	defer server.Close()

	resp, err := client.WebSearchWithCountry(context.Background(), "go programming", "JP")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "search", resp.Type)
	assert.NotNil(t, resp.Web)
}

// TestWebSearchWithLanguage tests the search with language helper function
func TestWebSearchWithLanguage(t *testing.T) {
	server, client := setupMockServer(t)
	defer server.Close()

	resp, err := client.WebSearchWithLanguage(context.Background(), "go programming", "ja")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "search", resp.Type)
	assert.NotNil(t, resp.Web)
}

// TestWebSearchNews tests the news search helper function
func TestWebSearchNews(t *testing.T) {
	server, client := setupMockServer(t)
	defer server.Close()

	resp, err := client.WebSearchNews(context.Background(), "go programming")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "search", resp.Type)
	assert.NotNil(t, resp.Web)
}

// TestWebSearchWithSafeSearch tests the safe search helper function
func TestWebSearchWithSafeSearch(t *testing.T) {
	server, client := setupMockServer(t)
	defer server.Close()

	resp, err := client.WebSearchWithSafeSearch(context.Background(), "go programming", SafeSearchStrict)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "search", resp.Type)
	assert.NotNil(t, resp.Web)
}

// TestWebSearchWithFreshness tests the freshness helper function
func TestWebSearchWithFreshness(t *testing.T) {
	server, client := setupMockServer(t)
	defer server.Close()

	resp, err := client.WebSearchWithFreshness(context.Background(), "go programming", FreshnessDay)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "search", resp.Type)
	assert.NotNil(t, resp.Web)
}

// TestWebSearchWithPagination tests the pagination helper function
func TestWebSearchWithPagination(t *testing.T) {
	server, client := setupMockServer(t)
	defer server.Close()

	resp, err := client.WebSearchWithPagination(context.Background(), "go programming", 5, 2)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "search", resp.Type)
	assert.NotNil(t, resp.Web)
}

// TestWebSearchResponseHelpers tests the helper methods on WebSearchResponse
func TestWebSearchResponseHelpers(t *testing.T) {
	// Load test data
	data, err := os.ReadFile("testdata/web_search_response.json")
	require.NoError(t, err)

	var response WebSearchResponse
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	// Test GetWebResults
	results := response.GetWebResults()
	assert.Len(t, results, 3)
	assert.Equal(t, "The Go Programming Language", results[0].Title)

	// Test HasMoreResults
	assert.True(t, response.HasMoreResults())

	// Test GetResultCount
	assert.Equal(t, 3, response.GetResultCount())

	// Test GetFirstResult
	firstResult := response.GetFirstResult()
	assert.NotNil(t, firstResult)
	assert.Equal(t, "The Go Programming Language", firstResult.Title)

	// Test IsWebResultEmpty
	assert.False(t, response.IsWebResultEmpty())

	// Test with nil web results
	emptyResponse := WebSearchResponse{}
	assert.Empty(t, emptyResponse.GetWebResults())
	assert.False(t, emptyResponse.HasMoreResults())
	assert.Equal(t, 0, emptyResponse.GetResultCount())
	assert.Nil(t, emptyResponse.GetFirstResult())
	assert.True(t, emptyResponse.IsWebResultEmpty())

	// Test with nil response
	var nilResponse *WebSearchResponse = nil
	assert.Empty(t, nilResponse.GetWebResults())
	assert.False(t, nilResponse.HasMoreResults())
	assert.Equal(t, 0, nilResponse.GetResultCount())
	assert.Nil(t, nilResponse.GetFirstResult())
	assert.True(t, nilResponse.IsWebResultEmpty())
}
