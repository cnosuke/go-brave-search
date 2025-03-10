# go-brave-search

[![Go Reference](https://pkg.go.dev/badge/github.com/cnosuke/go-brave-search.svg)](https://pkg.go.dev/github.com/cnosuke/go-brave-search)
[![Go Report Card](https://goreportcard.com/badge/github.com/cnosuke/go-brave-search)](https://goreportcard.com/report/github.com/cnosuke/go-brave-search)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client library for the Brave Search API. This library provides a simple and idiomatic Go interface to interact with Brave's search services.

## Features

- Simple, idiomatic Go API
- Support for Brave's Web Search API
- Configurable via functional options pattern
- Clear error handling
- Fully typed request and response structures

## Installation

```bash
go get github.com/cnosuke/go-brave-search
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	bravesearch "github.com/cnosuke/go-brave-search"
)

func main() {
	// Create a new client with your API key
	client, err := bravesearch.NewClient("your-api-key-here")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a context
	ctx := context.Background()

	// Perform a web search
	results, err := client.WebSearch(ctx, "brave search", nil)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// Print the results
	fmt.Printf("Search results for 'brave search':\n")
	for i, result := range results.Web.Results {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", result.URL)
		fmt.Printf("   Description: %s\n\n", result.Description)
	}
}
```

## Advanced Usage

### With Options

```go
package main

import (
	"context"
	"fmt"
	"log"

	bravesearch "github.com/cnosuke/go-brave-search"
)

func main() {
	// Create a client with options
	client, err := bravesearch.NewClient(
		"your-api-key-here",
		bravesearch.WithTimeout(30),
		bravesearch.WithDefaultCountry("JP"),
		bravesearch.WithDefaultSearchLanguage("ja"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a context
	ctx := context.Background()

	// Create search params
	params := &bravesearch.WebSearchParams{
		Count:      10,
		Offset:     0,
		SafeSearch: bravesearch.SafeSearchModerate,
		Freshness:  bravesearch.FreshnessWeek,
	}

	// Perform a web search with params
	results, err := client.WebSearch(ctx, "golang programming", params)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// Print the results
	fmt.Printf("Found %d results\n", len(results.Web.Results))
	for i, result := range results.Web.Results {
		fmt.Printf("%d. %s\n", i+1, result.Title)
	}
}
```

## API Reference

For detailed API documentation, see the [Go Reference](https://pkg.go.dev/github.com/cnosuke/go-brave-search).

### Client

`Client` is the main entry point for using the library:

```go
// Create a new client
client, err := bravesearch.NewClient("api-key")

// Create a client with options
client, err := bravesearch.NewClient(
    "api-key",
    bravesearch.WithTimeout(30),
    bravesearch.WithUserAgent("MyApp/1.0"),
)
```

### Web Search

```go
// Simple search
results, err := client.WebSearch(ctx, "query", nil)

// Search with parameters
params := &bravesearch.WebSearchParams{
    Count:       20,
    Offset:      0,
    Country:     "US",
    SearchLang:  "en",
    UILang:      "en-US",
    SafeSearch:  bravesearch.SafeSearchModerate,
    Freshness:   bravesearch.FreshnessMonth,
    Spellcheck:  true,
}
results, err := client.WebSearch(ctx, "query", params)
```

## Error Handling

The library provides detailed error information. Errors are wrapped with descriptive messages and can be unwrapped for more details.

```go
if err != nil {
    if bravesearch.IsRateLimitError(err) {
        // Handle rate limit error
    } else if bravesearch.IsAuthError(err) {
        // Handle authentication error
    } else {
        // Handle other errors
    }
}
```

## Configuration

The library supports several configuration options through functional options pattern:

```go
client, err := bravesearch.NewClient(
    "api-key",
    bravesearch.WithTimeout(30),                   // Request timeout in seconds
    bravesearch.WithRetries(3),                    // Number of retries on transient errors
    bravesearch.WithUserAgent("MyApp/1.0"),        // Custom User-Agent
    bravesearch.WithDefaultCountry("JP"),          // Default country for searches
    bravesearch.WithDefaultSearchLanguage("ja"),   // Default search language
    bravesearch.WithDefaultUILanguage("ja-JP"),    // Default UI language
)
```

## Development Status

This library is currently in active development. While it's functional and tested, we're continuously improving it. Feedback and contributions are welcome!

## License

MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- This library is not officially associated with or endorsed by Brave Software, Inc.
- Thanks to Brave for providing an excellent search API
