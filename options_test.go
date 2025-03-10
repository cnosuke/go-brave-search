package bravesearch

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestWithTimeout tests the WithTimeout option
func TestWithTimeout(t *testing.T) {
	config := &ClientConfig{}
	
	// Test with valid timeout
	err := WithTimeout(30)(config)
	assert.NoError(t, err)
	assert.Equal(t, 30*time.Second, config.Timeout)
	
	// Test with zero timeout
	err = WithTimeout(0)(config)
	assert.NoError(t, err)
	assert.Equal(t, 0*time.Second, config.Timeout)
	
	// Test with negative timeout (still sets the value, though it's not practical)
	err = WithTimeout(-10)(config)
	assert.NoError(t, err)
	assert.Equal(t, -10*time.Second, config.Timeout)
}

// TestWithRetries tests the WithRetries option
func TestWithRetries(t *testing.T) {
	config := &ClientConfig{}
	
	// Test with valid retries
	err := WithRetries(3)(config)
	assert.NoError(t, err)
	assert.Equal(t, 3, config.MaxRetries)
	
	// Test with zero retries
	err = WithRetries(0)(config)
	assert.NoError(t, err)
	assert.Equal(t, 0, config.MaxRetries)
	
	// Test with negative retries (should error)
	err = WithRetries(-1)(config)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidParameters, err)
}

// TestWithUserAgent tests the WithUserAgent option
func TestWithUserAgent(t *testing.T) {
	config := &ClientConfig{}
	
	// Test with valid user agent
	err := WithUserAgent("custom-agent")(config)
	assert.NoError(t, err)
	assert.Equal(t, "custom-agent", config.UserAgent)
	
	// Test with empty user agent
	err = WithUserAgent("")(config)
	assert.NoError(t, err)
	assert.Equal(t, "", config.UserAgent)
}

// TestWithBaseURL tests the WithBaseURL option
func TestWithBaseURL(t *testing.T) {
	config := &ClientConfig{}
	
	// Test with valid base URL
	err := WithBaseURL("https://custom-api.example.com")(config)
	assert.NoError(t, err)
	assert.Equal(t, "https://custom-api.example.com", config.BaseURL)
	
	// Test with empty base URL
	err = WithBaseURL("")(config)
	assert.NoError(t, err)
	assert.Equal(t, "", config.BaseURL)
}

// TestWithHTTPClient tests the WithHTTPClient option
func TestWithHTTPClient(t *testing.T) {
	config := &ClientConfig{}
	customClient := &http.Client{Timeout: 60 * time.Second}
	
	// Test with custom HTTP client
	err := WithHTTPClient(customClient)(config)
	assert.NoError(t, err)
	assert.Equal(t, customClient, config.HTTPClient)
	
	// Test with nil HTTP client
	err = WithHTTPClient(nil)(config)
	assert.NoError(t, err)
	assert.Nil(t, config.HTTPClient)
}

// TestWithDefaultCountry tests the WithDefaultCountry option
func TestWithDefaultCountry(t *testing.T) {
	config := &ClientConfig{}
	
	// Test with valid country
	err := WithDefaultCountry("JP")(config)
	assert.NoError(t, err)
	assert.Equal(t, "JP", config.DefaultCountry)
	
	// Test with empty country
	err = WithDefaultCountry("")(config)
	assert.NoError(t, err)
	assert.Equal(t, "", config.DefaultCountry)
}

// TestWithDefaultSearchLanguage tests the WithDefaultSearchLanguage option
func TestWithDefaultSearchLanguage(t *testing.T) {
	config := &ClientConfig{}
	
	// Test with valid language
	err := WithDefaultSearchLanguage("ja")(config)
	assert.NoError(t, err)
	assert.Equal(t, "ja", config.DefaultSearchLang)
	
	// Test with empty language
	err = WithDefaultSearchLanguage("")(config)
	assert.NoError(t, err)
	assert.Equal(t, "", config.DefaultSearchLang)
}

// TestWithDefaultUILanguage tests the WithDefaultUILanguage option
func TestWithDefaultUILanguage(t *testing.T) {
	config := &ClientConfig{}
	
	// Test with valid UI language
	err := WithDefaultUILanguage("ja-JP")(config)
	assert.NoError(t, err)
	assert.Equal(t, "ja-JP", config.DefaultUILang)
	
	// Test with empty UI language
	err = WithDefaultUILanguage("")(config)
	assert.NoError(t, err)
	assert.Equal(t, "", config.DefaultUILang)
}

// TestApplyOptions tests applying multiple options
func TestApplyOptions(t *testing.T) {
	config := &ClientConfig{}
	
	// Define options
	options := []ClientOption{
		WithTimeout(45),
		WithRetries(5),
		WithUserAgent("test-agent"),
		WithBaseURL("https://test.example.com"),
		WithDefaultCountry("FR"),
		WithDefaultSearchLanguage("fr"),
		WithDefaultUILanguage("fr-FR"),
	}
	
	// Apply options
	err := applyOptions(config, options...)
	assert.NoError(t, err)
	
	// Verify all options were applied
	assert.Equal(t, 45*time.Second, config.Timeout)
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, "test-agent", config.UserAgent)
	assert.Equal(t, "https://test.example.com", config.BaseURL)
	assert.Equal(t, "FR", config.DefaultCountry)
	assert.Equal(t, "fr", config.DefaultSearchLang)
	assert.Equal(t, "fr-FR", config.DefaultUILang)
}

// TestApplyOptionsWithError tests error handling in apply options
func TestApplyOptionsWithError(t *testing.T) {
	config := &ClientConfig{}
	
	// Define options with one that will error
	options := []ClientOption{
		WithTimeout(45),
		WithRetries(-1), // This should cause an error
		WithUserAgent("test-agent"),
	}
	
	// Apply options
	err := applyOptions(config, options...)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidParameters, err)
	
	// Verify first option was applied but not the third
	assert.Equal(t, 45*time.Second, config.Timeout)
	assert.NotEqual(t, "test-agent", config.UserAgent)
}
