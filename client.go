package bravesearch

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client is the API client for Brave Search
type Client struct {
	config ClientConfig
	http   *http.Client
}

// NewClient creates a new Brave Search API client
func NewClient(apiKey string, options ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	// Default configuration
	config := ClientConfig{
		APIKey:            apiKey,
		BaseURL:           BaseURL,
		Timeout:           time.Duration(DefaultTimeout) * time.Second,
		MaxRetries:        DefaultMaxRetries,
		UserAgent:         DefaultUserAgent,
		DefaultCountry:    DefaultCountry,
		DefaultSearchLang: DefaultSearchLang,
		DefaultUILang:     DefaultUILang,
	}

	// Apply options
	if err := applyOptions(&config, options...); err != nil {
		return nil, err
	}

	// Create HTTP client if not provided
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	client := &Client{
		config: config,
		http:   httpClient,
	}

	return client, nil
}

// WebSearch performs a web search
func (c *Client) WebSearch(ctx context.Context, query string, params *WebSearchParams) (*WebSearchResponse, error) {
	if query == "" {
		return nil, ErrEmptyQuery
	}

	if len(query) > 400 || len(strings.Fields(query)) > 50 {
		return nil, ErrQueryTooLong
	}

	// Create a copy of params or initialize a new one
	searchParams := &WebSearchParams{}
	if params != nil {
		*searchParams = *params
	}

	// Set query
	searchParams.Query = query

	// Apply defaults if not set
	if searchParams.Country == "" {
		searchParams.Country = c.config.DefaultCountry
	}
	if searchParams.SearchLang == "" {
		searchParams.SearchLang = c.config.DefaultSearchLang
	}
	if searchParams.UILang == "" {
		searchParams.UILang = c.config.DefaultUILang
	}
	if searchParams.Count == 0 {
		searchParams.Count = DefaultCount
	}
	if searchParams.SafeSearch == "" {
		searchParams.SafeSearch = DefaultSafeSearch
	}

	// Build URL
	requestURL, err := c.buildRequestURL(WebSearchEndpoint, searchParams)
	if err != nil {
		return nil, err
	}

	// Make the request
	var response WebSearchResponse
	if err := c.makeRequest(ctx, http.MethodGet, requestURL, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// buildRequestURL builds the request URL with query parameters
func (c *Client) buildRequestURL(endpoint string, params *WebSearchParams) (string, error) {
	// Ensure baseURL ends with slash if endpoint doesn't start with one
	baseURL := c.config.BaseURL
	if !strings.HasSuffix(baseURL, "/") && !strings.HasPrefix(endpoint, "/") {
		baseURL += "/"
	}
	baseURL += endpoint

	// Build query string
	values := url.Values{}
	if params.Query != "" {
		values.Add("q", params.Query)
	}
	if params.Country != "" {
		values.Add("country", params.Country)
	}
	if params.SearchLang != "" {
		values.Add("search_lang", params.SearchLang)
	}
	if params.UILang != "" {
		values.Add("ui_lang", params.UILang)
	}
	if params.Count > 0 {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Offset > 0 {
		values.Add("offset", strconv.Itoa(params.Offset))
	}
	if params.SafeSearch != "" {
		values.Add("safesearch", params.SafeSearch)
	}
	if params.Freshness != "" {
		values.Add("freshness", params.Freshness)
	}
	values.Add("text_decorations", strconv.FormatBool(params.TextDecorations))
	values.Add("spellcheck", strconv.FormatBool(params.Spellcheck))
	if params.ResultFilter != "" {
		values.Add("result_filter", params.ResultFilter)
	}
	if params.Goggles != "" {
		values.Add("goggles", params.Goggles)
	}
	if params.Units != "" {
		values.Add("units", params.Units)
	}
	if params.ExtraSnippets {
		values.Add("extra_snippets", "true")
	}
	if params.Summary {
		values.Add("summary", "true")
	}

	// Append query string to URL
	return baseURL + "?" + values.Encode(), nil
}

// makeRequest makes an HTTP request to the API
func (c *Client) makeRequest(ctx context.Context, method, url string, body interface{}, result interface{}) error {
	var bodyReader io.Reader

	// Prepare request body if any
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set(HeaderAccept, MIMETypeJSON)
	req.Header.Set(HeaderAcceptEncoding, MIMETypeGzip)
	req.Header.Set(HeaderUserAgent, c.config.UserAgent)
	req.Header.Set(HeaderSubscriptionToken, c.config.APIKey)
	req.Header.Set(HeaderCacheControl, "no-cache")

	if body != nil {
		req.Header.Set("Content-Type", MIMETypeJSON)
	}

	// Make the request with retries
	var resp *http.Response
	var respErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		resp, respErr = c.http.Do(req)
		if respErr == nil && resp.StatusCode < 500 {
			// Success or non-retriable error
			break
		}

		// If this was the last attempt, return the error
		if attempt == c.config.MaxRetries {
			if respErr != nil {
				return respErr
			}
			return NewHTTPError(resp)
		}

		// Close response body if any
		if resp != nil {
			resp.Body.Close()
		}

		// Add exponential backoff
		backoffTime := time.Duration(1<<uint(attempt)) * 100 * time.Millisecond
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoffTime):
			// Continue with retry
		}
	}

	// Handle HTTP errors
	if resp == nil {
		return fmt.Errorf("no response: %w", respErr)
	}
	defer resp.Body.Close()

	// Handle HTTP error status codes
	if resp.StatusCode != http.StatusOK {
		var bodyReader io.ReadCloser
		if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
			bodyReader, _ = gzip.NewReader(resp.Body)
			defer bodyReader.Close()

			respBody, _ := io.ReadAll(bodyReader)
			resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
		} else {
			// For debugging, print the response body if there's an error
			respBody, _ := io.ReadAll(resp.Body)
			// Create a new response with the same body for further processing
			resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
		}

		return NewHTTPError(resp)
	}

	// Parse rate limit headers
	c.parseRateLimitHeaders(resp)

	// Parse response body
	if result != nil {
		// Check if content is gzipped
		var bodyReader io.ReadCloser
		contentEncoding := resp.Header.Get("Content-Encoding")

		if strings.Contains(contentEncoding, "gzip") {
			// Import required for gzip
			var err error
			bodyReader, err = gzip.NewReader(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to create gzip reader: %w", err)
			}
			defer bodyReader.Close()
		} else {
			bodyReader = resp.Body
		}

		// Read the body
		body, err := io.ReadAll(bodyReader)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, result); err != nil {
			return &APIError{
				StatusCode: resp.StatusCode,
				Message:    "Failed to parse response",
				Err:        ErrInvalidResponse,
			}
		}
	}

	return nil
}

// parseRateLimitHeaders parses rate limit information from response headers
func (c *Client) parseRateLimitHeaders(resp *http.Response) *RateLimit {
	rateLimit := &RateLimit{}

	// Parse X-RateLimit-Limit
	if limitStr := resp.Header.Get(HeaderRateLimitLimit); limitStr != "" {
		limits := strings.Split(limitStr, ",")
		if len(limits) > 0 {
			limit, err := strconv.Atoi(strings.TrimSpace(limits[0]))
			if err == nil {
				rateLimit.Limit = limit
			}
		}
	}

	// Parse X-RateLimit-Remaining
	if remainingStr := resp.Header.Get(HeaderRateLimitRemaining); remainingStr != "" {
		remaining := strings.Split(remainingStr, ",")
		if len(remaining) > 0 {
			remain, err := strconv.Atoi(strings.TrimSpace(remaining[0]))
			if err == nil {
				rateLimit.Remaining = remain
			}
		}
	}

	// Parse X-RateLimit-Reset
	if resetStr := resp.Header.Get(HeaderRateLimitReset); resetStr != "" {
		resets := strings.Split(resetStr, ",")
		if len(resets) > 0 {
			reset, err := strconv.Atoi(strings.TrimSpace(resets[0]))
			if err == nil {
				rateLimit.Reset = reset
			}
		}
	}

	return rateLimit
}
