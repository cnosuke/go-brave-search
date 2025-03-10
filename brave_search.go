// Package bravesearch provides a Go client for the Brave Search API.
//
// This package allows Go applications to interact with Brave's search services,
// including the Web Search API. It provides a simple, idiomatic interface
// with support for configuration via functional options.
//
// Example usage:
//
//	client, err := bravesearch.NewClient("your-api-key")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	results, err := client.WebSearch(context.Background(), "brave search", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, result := range results.Web.Results {
//		fmt.Printf("%s - %s\n", result.Title, result.URL)
//	}
package bravesearch

// Version is the current version of the client library
const Version = "0.1.0"

// UserAgentPrefix is the prefix for the user agent string
const UserAgentPrefix = "go-brave-search"

// GetVersion returns the current version of the client library
func GetVersion() string {
	return Version
}

// GetUserAgent returns the user agent string for this library
func GetUserAgent() string {
	return UserAgentPrefix + "/" + Version
}
