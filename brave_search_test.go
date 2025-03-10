package bravesearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetVersion tests the GetVersion function
func TestGetVersion(t *testing.T) {
	version := GetVersion()
	assert.Equal(t, Version, version)
	assert.NotEmpty(t, version)
}

// TestGetUserAgent tests the GetUserAgent function
func TestGetUserAgent(t *testing.T) {
	userAgent := GetUserAgent()
	assert.Equal(t, UserAgentPrefix+"/"+Version, userAgent)
	assert.Contains(t, userAgent, "go-brave-search")
	assert.Contains(t, userAgent, Version)
}
