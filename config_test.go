package bravesearch

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewDefaultConfig tests creating a default configuration
func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()
	assert.NotNil(t, config)
	
	// Check that defaults are set correctly
	assert.Equal(t, BaseURL, config.BaseURL)
	assert.Equal(t, time.Duration(DefaultTimeout)*time.Second, config.Timeout)
	assert.Equal(t, DefaultMaxRetries, config.MaxRetries)
	assert.Equal(t, DefaultUserAgent, config.UserAgent)
	assert.Equal(t, DefaultCountry, config.DefaultCountry)
	assert.Equal(t, DefaultSearchLang, config.DefaultSearchLang)
	assert.Equal(t, DefaultUILang, config.DefaultUILang)
}

// TestWithConfig tests creating a client option with configuration
func TestWithConfig(t *testing.T) {
	// Create a custom config
	customConfig := &Config{
		APIKey:           "custom-api-key",
		BaseURL:          "https://custom.example.com",
		HTTPClient:       &http.Client{Timeout: 60 * time.Second},
		Timeout:          45 * time.Second,
		MaxRetries:       5,
		UserAgent:        "custom-agent",
		DefaultCountry:   "FR",
		DefaultSearchLang: "fr",
		DefaultUILang:    "fr-FR",
	}
	
	// Create client config
	clientConfig := &ClientConfig{}
	
	// Apply WithConfig option
	err := WithConfig(customConfig)(clientConfig)
	assert.NoError(t, err)
	
	// Verify config was applied
	assert.Equal(t, "custom-api-key", clientConfig.APIKey)
	assert.Equal(t, "https://custom.example.com", clientConfig.BaseURL)
	assert.Equal(t, customConfig.HTTPClient, clientConfig.HTTPClient)
	assert.Equal(t, 45*time.Second, clientConfig.Timeout)
	assert.Equal(t, 5, clientConfig.MaxRetries)
	assert.Equal(t, "custom-agent", clientConfig.UserAgent)
	assert.Equal(t, "FR", clientConfig.DefaultCountry)
	assert.Equal(t, "fr", clientConfig.DefaultSearchLang)
	assert.Equal(t, "fr-FR", clientConfig.DefaultUILang)
	
	// Test with nil config (should do nothing)
	clientConfig = &ClientConfig{
		APIKey: "original-key",
	}
	err = WithConfig(nil)(clientConfig)
	assert.NoError(t, err)
	assert.Equal(t, "original-key", clientConfig.APIKey)
	
	// Test with partial config (should only update specified fields)
	partialConfig := &Config{
		APIKey:     "new-api-key",
		MaxRetries: 3,
	}
	clientConfig = &ClientConfig{
		APIKey:     "original-key",
		BaseURL:    "original-url",
		MaxRetries: 1,
	}
	err = WithConfig(partialConfig)(clientConfig)
	assert.NoError(t, err)
	assert.Equal(t, "new-api-key", clientConfig.APIKey)
	assert.Equal(t, "original-url", clientConfig.BaseURL) // Unchanged
	assert.Equal(t, 3, clientConfig.MaxRetries)
}

// TestValidateConfig tests configuration validation
func TestValidateConfig(t *testing.T) {
	// Test valid config
	validConfig := &ClientConfig{
		APIKey: "test-key",
	}
	err := ValidateConfig(validConfig)
	assert.NoError(t, err)
	
	// Test missing API key
	invalidConfig := &ClientConfig{
		APIKey: "",
	}
	err = ValidateConfig(invalidConfig)
	assert.Error(t, err)
	assert.Equal(t, ErrMissingAPIKey, err)
}
