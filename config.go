package bravesearch

import (
	"net/http"
	"time"
)

// Config provides configuration for the Brave Search API client.
type Config struct {
	// APIKey is the authentication token for the API
	APIKey string

	// BaseURL is the base URL for API requests
	BaseURL string

	// HTTPClient is the HTTP client to use for requests
	HTTPClient *http.Client

	// Timeout is the timeout for requests
	Timeout time.Duration

	// MaxRetries is the maximum number of retries for failed requests
	MaxRetries int

	// UserAgent is the User-Agent header to use for requests
	UserAgent string

	// DefaultCountry is the default country code for searches
	DefaultCountry string

	// DefaultSearchLang is the default search language
	DefaultSearchLang string

	// DefaultUILang is the default UI language
	DefaultUILang string
}

// NewDefaultConfig creates a new default configuration
func NewDefaultConfig() *Config {
	return &Config{
		BaseURL:          BaseURL,
		Timeout:          time.Duration(DefaultTimeout) * time.Second,
		MaxRetries:       DefaultMaxRetries,
		UserAgent:        DefaultUserAgent,
		DefaultCountry:   DefaultCountry,
		DefaultSearchLang: DefaultSearchLang,
		DefaultUILang:    DefaultUILang,
	}
}

// WithConfig creates a new client with the given configuration
func WithConfig(config *Config) ClientOption {
	return func(c *ClientConfig) error {
		if config == nil {
			return nil
		}

		if config.APIKey != "" {
			c.APIKey = config.APIKey
		}
		if config.BaseURL != "" {
			c.BaseURL = config.BaseURL
		}
		if config.HTTPClient != nil {
			c.HTTPClient = config.HTTPClient
		}
		if config.Timeout > 0 {
			c.Timeout = config.Timeout
		}
		if config.MaxRetries >= 0 {
			c.MaxRetries = config.MaxRetries
		}
		if config.UserAgent != "" {
			c.UserAgent = config.UserAgent
		}
		if config.DefaultCountry != "" {
			c.DefaultCountry = config.DefaultCountry
		}
		if config.DefaultSearchLang != "" {
			c.DefaultSearchLang = config.DefaultSearchLang
		}
		if config.DefaultUILang != "" {
			c.DefaultUILang = config.DefaultUILang
		}

		return nil
	}
}

// ValidateConfig validates the client configuration
func ValidateConfig(config *ClientConfig) error {
	if config.APIKey == "" {
		return ErrMissingAPIKey
	}
	
	return nil
}
