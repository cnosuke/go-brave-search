package bravesearch

import (
	"net/http"
	"time"
)

// ClientOption is a function that can be used to configure a Client
type ClientOption func(*ClientConfig) error

// WithTimeout sets the timeout for requests
func WithTimeout(seconds int) ClientOption {
	return func(c *ClientConfig) error {
		c.Timeout = time.Duration(seconds) * time.Second
		return nil
	}
}

// WithRetries sets the maximum number of retries for requests
func WithRetries(retries int) ClientOption {
	return func(c *ClientConfig) error {
		if retries < 0 {
			return ErrInvalidParameters
		}
		c.MaxRetries = retries
		return nil
	}
}

// WithUserAgent sets the User-Agent header for requests
func WithUserAgent(userAgent string) ClientOption {
	return func(c *ClientConfig) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithBaseURL sets the base URL for the API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *ClientConfig) error {
		c.BaseURL = baseURL
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *ClientConfig) error {
		c.HTTPClient = client
		return nil
	}
}

// WithDefaultCountry sets the default country for requests
func WithDefaultCountry(country string) ClientOption {
	return func(c *ClientConfig) error {
		c.DefaultCountry = country
		return nil
	}
}

// WithDefaultSearchLanguage sets the default search language for requests
func WithDefaultSearchLanguage(lang string) ClientOption {
	return func(c *ClientConfig) error {
		c.DefaultSearchLang = lang
		return nil
	}
}

// WithDefaultUILanguage sets the default UI language for requests
func WithDefaultUILanguage(lang string) ClientOption {
	return func(c *ClientConfig) error {
		c.DefaultUILang = lang
		return nil
	}
}

// applyOptions applies the given options to the config
func applyOptions(config *ClientConfig, options ...ClientOption) error {
	for _, option := range options {
		if err := option(config); err != nil {
			return err
		}
	}
	return nil
}
