package bravesearch

import (
	"net/http"
	"time"
)

// ClientConfig holds the configuration for the API client
type ClientConfig struct {
	APIKey           string
	BaseURL          string
	Timeout          time.Duration
	MaxRetries       int
	UserAgent        string
	DefaultCountry   string
	DefaultSearchLang string
	DefaultUILang    string
	HTTPClient       *http.Client
}

// WebSearchParams holds the parameters for a web search request
type WebSearchParams struct {
	// Required parameters
	Query string `url:"q,omitempty"`

	// Optional parameters
	Country         string `url:"country,omitempty"`
	SearchLang      string `url:"search_lang,omitempty"`
	UILang          string `url:"ui_lang,omitempty"`
	Count           int    `url:"count,omitempty"`
	Offset          int    `url:"offset,omitempty"`
	SafeSearch      string `url:"safesearch,omitempty"`
	Freshness       string `url:"freshness,omitempty"`
	TextDecorations bool   `url:"text_decorations,omitempty"`
	Spellcheck      bool   `url:"spellcheck,omitempty"`
	ResultFilter    string `url:"result_filter,omitempty"`
	Goggles         string `url:"goggles,omitempty"`
	Units           string `url:"units,omitempty"`
	ExtraSnippets   bool   `url:"extra_snippets,omitempty"`
	Summary         bool   `url:"summary,omitempty"`
}

// WebSearchResponse represents the top-level response from the Web Search API
type WebSearchResponse struct {
	Type        string          `json:"type"`
	Discussions *Discussions    `json:"discussions,omitempty"`
	FAQ         *FAQ            `json:"faq,omitempty"`
	Infobox     *GraphInfobox   `json:"infobox,omitempty"`
	Locations   *Locations      `json:"locations,omitempty"`
	Mixed       *MixedResponse  `json:"mixed,omitempty"`
	News        *News           `json:"news,omitempty"`
	Query       *Query          `json:"query,omitempty"`
	Videos      *Videos         `json:"videos,omitempty"`
	Web         *Search         `json:"web,omitempty"`
	Summarizer  *Summarizer     `json:"summarizer,omitempty"`
}

// Search represents a collection of web search results
type Search struct {
	Type           string         `json:"type"`
	Results        []SearchResult `json:"results"`
	FamilyFriendly bool           `json:"family_friendly"`
}

// SearchResult represents an individual web search result
type SearchResult struct {
	Title          string   `json:"title"`
	URL            string   `json:"url"`
	IsSourceLocal  bool     `json:"is_source_local"`
	IsSourceBoth   bool     `json:"is_source_both"`
	Description    string   `json:"description,omitempty"`
	PageAge        string   `json:"page_age,omitempty"`
	PageFetched    string   `json:"page_fetched,omitempty"`
	Profile        *Profile `json:"profile,omitempty"`
	Language       string   `json:"language,omitempty"`
	FamilyFriendly bool     `json:"family_friendly"`
	Type           string   `json:"type"`
	Subtype        string   `json:"subtype,omitempty"`
	IsLive         bool     `json:"is_live,omitempty"`
	DeepResults    *DeepResults `json:"deep_results,omitempty"`
	MetaURL        *MetaURL     `json:"meta_url,omitempty"`
	Thumbnail      *Thumbnail   `json:"thumbnail,omitempty"`
	Age            string       `json:"age,omitempty"`
}

// Profile represents profile information associated with a search result
type Profile struct {
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	LongName string `json:"long_name,omitempty"`
	Img      string `json:"img,omitempty"`
}

// DeepResults represents additional links or features for a search result
type DeepResults struct {
	Buttons []ButtonResult `json:"buttons,omitempty"`
}

// ButtonResult represents a button in deep results
type ButtonResult struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// MetaURL represents metadata about the URL
type MetaURL struct {
	Scheme   string `json:"scheme,omitempty"`
	Netloc   string `json:"netloc,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Favicon  string `json:"favicon,omitempty"`
	Path     string `json:"path,omitempty"`
}

// Thumbnail represents an image thumbnail
type Thumbnail struct {
	Src      string `json:"src,omitempty"`
	Original string `json:"original,omitempty"`
	Logo     bool   `json:"logo,omitempty"`
}

// Query represents query information
type Query struct {
	Original           string `json:"original"`
	Altered            string `json:"altered,omitempty"`
	ShowStrictWarning  bool   `json:"show_strict_warning"`
	IsNavigational     bool   `json:"is_navigational"`
	IsNewsBreaking     bool   `json:"is_news_breaking"`
	SpellcheckOff      bool   `json:"spellcheck_off"`
	Country            string `json:"country"`
	BadResults         bool   `json:"bad_results"`
	ShouldFallback     bool   `json:"should_fallback"`
	PostalCode         string `json:"postal_code"`
	City               string `json:"city"`
	HeaderCountry      string `json:"header_country"`
	MoreResultsAvailable bool   `json:"more_results_available"`
	State              string `json:"state"`
}

// MixedResponse represents a collection of mixed result types
type MixedResponse struct {
	Type string           `json:"type"`
	Main []MixedResultRef `json:"main"`
	Top  []MixedResultRef `json:"top"`
	Side []MixedResultRef `json:"side"`
}

// MixedResultRef references a result in another section
type MixedResultRef struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	All   bool   `json:"all"`
}

// The following types are included for completeness, but aren't fully implemented
// since we're focusing on Web results only

// Discussions represents forum discussions
type Discussions struct {
	Type    string `json:"type"`
	Results []any  `json:"results,omitempty"`
}

// FAQ represents frequently asked questions
type FAQ struct {
	Type    string `json:"type"`
	Results []any  `json:"results,omitempty"`
}

// GraphInfobox represents an infobox
type GraphInfobox struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

// Locations represents location results
type Locations struct {
	Type    string `json:"type"`
	Results []any  `json:"results,omitempty"`
}

// News represents news results
type News struct {
	Type    string `json:"type"`
	Results []any  `json:"results,omitempty"`
}

// Videos represents video results
type Videos struct {
	Type    string `json:"type"`
	Results []any  `json:"results,omitempty"`
}

// Summarizer represents summary results
type Summarizer struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

// RateLimit represents rate limit information
type RateLimit struct {
	Limit     int
	Remaining int
	Reset     int
}
