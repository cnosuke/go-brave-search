package bravesearch

// API Endpoints
const (
	// BaseURL is the base URL for Brave Search API
	BaseURL = "https://api.search.brave.com/res/v1"

	// WebSearchEndpoint is the endpoint for web search
	WebSearchEndpoint = "/web/search"
)

// SafeSearch options
const (
	SafeSearchOff      = "off"
	SafeSearchModerate = "moderate"
	SafeSearchStrict   = "strict"
)

// Freshness options
const (
	FreshnessDay   = "pd"  // Past day
	FreshnessWeek  = "pw"  // Past week
	FreshnessMonth = "pm"  // Past month
	FreshnessYear  = "py"  // Past year
)

// Default values
const (
	DefaultCount        = 20
	DefaultOffset       = 0
	DefaultSafeSearch   = SafeSearchModerate
	DefaultCountry      = "US"
	DefaultSearchLang   = "en"
	DefaultUILang       = "en-US"
	DefaultTimeout      = 30 // seconds
	DefaultMaxRetries   = 2
	DefaultUserAgent    = "go-brave-search/1.0"
	DefaultTextDecor    = true
	DefaultSpellCheck   = true
)

// HTTP Headers
const (
	HeaderAccept             = "Accept"
	HeaderAcceptEncoding     = "Accept-Encoding"
	HeaderUserAgent          = "User-Agent"
	HeaderSubscriptionToken  = "X-Subscription-Token"
	HeaderCacheControl       = "Cache-Control"
	HeaderLocLatitude        = "X-Loc-Lat"
	HeaderLocLongitude       = "X-Loc-Long"
	HeaderLocTimezone        = "X-Loc-Timezone"
	HeaderLocCity            = "X-Loc-City"
	HeaderLocState           = "X-Loc-State"
	HeaderLocStateName       = "X-Loc-State-Name"
	HeaderLocCountry         = "X-Loc-Country"
	HeaderLocPostalCode      = "X-Loc-Postal-Code"
)

// Response Headers
const (
	HeaderRateLimitLimit     = "X-RateLimit-Limit"
	HeaderRateLimitPolicy    = "X-RateLimit-Policy"
	HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRateLimitReset     = "X-RateLimit-Reset"
)

// MIME types
const (
	MIMETypeJSON           = "application/json"
	MIMETypeGzip           = "gzip"
)

// Result filters
const (
	ResultFilterDiscussions = "discussions"
	ResultFilterFaq         = "faq"
	ResultFilterInfobox     = "infobox"
	ResultFilterNews        = "news"
	ResultFilterQuery       = "query"
	ResultFilterSummarizer  = "summarizer"
	ResultFilterVideos      = "videos"
	ResultFilterWeb         = "web"
	ResultFilterLocations   = "locations"
)

// Units
const (
	UnitMetric   = "metric"
	UnitImperial = "imperial"
)
