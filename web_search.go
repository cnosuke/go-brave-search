package bravesearch

import (
	"context"
)

// NewWebSearchParams creates a new WebSearchParams with default values
func NewWebSearchParams() *WebSearchParams {
	return &WebSearchParams{
		Count:           DefaultCount,
		Offset:          DefaultOffset,
		SafeSearch:      DefaultSafeSearch,
		TextDecorations: DefaultTextDecor,
		Spellcheck:      DefaultSpellCheck,
	}
}

// WebSearchWithCountry performs a web search with a specific country
func (c *Client) WebSearchWithCountry(ctx context.Context, query string, country string) (*WebSearchResponse, error) {
	params := NewWebSearchParams()
	params.Country = country
	return c.WebSearch(ctx, query, params)
}

// WebSearchWithLanguage performs a web search with a specific search language
func (c *Client) WebSearchWithLanguage(ctx context.Context, query string, lang string) (*WebSearchResponse, error) {
	params := NewWebSearchParams()
	params.SearchLang = lang
	return c.WebSearch(ctx, query, params)
}

// WebSearchNews performs a web search filtered to news results
func (c *Client) WebSearchNews(ctx context.Context, query string) (*WebSearchResponse, error) {
	params := NewWebSearchParams()
	params.ResultFilter = ResultFilterNews
	return c.WebSearch(ctx, query, params)
}

// WebSearchVideos performs a web search filtered to video results
func (c *Client) WebSearchVideos(ctx context.Context, query string) (*WebSearchResponse, error) {
	params := NewWebSearchParams()
	params.ResultFilter = ResultFilterVideos
	return c.WebSearch(ctx, query, params)
}

// WebSearchWithSafeSearch performs a web search with a specific SafeSearch setting
func (c *Client) WebSearchWithSafeSearch(ctx context.Context, query string, safeSearch string) (*WebSearchResponse, error) {
	params := NewWebSearchParams()
	params.SafeSearch = safeSearch
	return c.WebSearch(ctx, query, params)
}

// WebSearchWithFreshness performs a web search with a specific freshness setting
func (c *Client) WebSearchWithFreshness(ctx context.Context, query string, freshness string) (*WebSearchResponse, error) {
	params := NewWebSearchParams()
	params.Freshness = freshness
	return c.WebSearch(ctx, query, params)
}

// WebSearchWithPagination performs a web search with pagination
func (c *Client) WebSearchWithPagination(ctx context.Context, query string, count, offset int) (*WebSearchResponse, error) {
	params := NewWebSearchParams()
	params.Count = count
	params.Offset = offset
	return c.WebSearch(ctx, query, params)
}

// WebSearchSummary performs a web search with summary enabled
func (c *Client) WebSearchSummary(ctx context.Context, query string) (*WebSearchResponse, error) {
	params := NewWebSearchParams()
	params.Summary = true
	return c.WebSearch(ctx, query, params)
}

// WebSearchWithUnits performs a web search with specific measurement units
func (c *Client) WebSearchWithUnits(ctx context.Context, query string, units string) (*WebSearchResponse, error) {
	params := NewWebSearchParams()
	params.Units = units
	return c.WebSearch(ctx, query, params)
}

// WebSearchRecent performs a web search for recent content
func (c *Client) WebSearchRecent(ctx context.Context, query string) (*WebSearchResponse, error) {
	params := NewWebSearchParams()
	params.Freshness = FreshnessDay
	return c.WebSearch(ctx, query, params)
}

// GetWebResults is a helper function to extract web results from the response
func (r *WebSearchResponse) GetWebResults() []SearchResult {
	if r == nil || r.Web == nil {
		return []SearchResult{}
	}
	return r.Web.Results
}

// HasMoreResults checks if the search has more results available
func (r *WebSearchResponse) HasMoreResults() bool {
	if r == nil || r.Query == nil {
		return false
	}
	return r.Query.MoreResultsAvailable
}

// GetResultCount returns the number of web results
func (r *WebSearchResponse) GetResultCount() int {
	if r == nil || r.Web == nil {
		return 0
	}
	return len(r.Web.Results)
}

// GetFirstResult returns the first web result or nil if no results
func (r *WebSearchResponse) GetFirstResult() *SearchResult {
	if r == nil || r.Web == nil || len(r.Web.Results) == 0 {
		return nil
	}
	return &r.Web.Results[0]
}

// IsWebResultEmpty checks if the web results are empty
func (r *WebSearchResponse) IsWebResultEmpty() bool {
	return r == nil || r.Web == nil || len(r.Web.Results) == 0
}
