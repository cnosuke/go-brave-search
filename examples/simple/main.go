package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	bravesearch "github.com/cnosuke/go-brave-search"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("BRAVE_API_KEY")
	if apiKey == "" {
		log.Fatal("BRAVE_API_KEY environment variable is required")
	}

	// Create a new client with options
	// Note: search_lang should be 'jp' for Japanese according to Brave API example
	client, err := bravesearch.NewClient(
		apiKey,
		bravesearch.WithTimeout(30),
		bravesearch.WithDefaultCountry("JP"),
		bravesearch.WithDefaultSearchLanguage("jp"), // Using 'jp' instead of 'ja'
		bravesearch.WithDefaultUILanguage("ja-JP"),  // Adding UI language
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get search query from command line arguments or use default
	query := "Go programming language"
	if len(os.Args) > 1 {
		query = os.Args[1]
	}

	fmt.Printf("Searching for: %s\n\n", query)

	// Simple search with default parameters
	results, err := client.WebSearch(ctx, query, nil)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// Print the results
	fmt.Printf("Found %d results\n\n", len(results.Web.Results))
	for i, result := range results.Web.Results {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", result.URL)
		fmt.Printf("   Description: %s\n", result.Description)
		if result.Profile != nil {
			fmt.Printf("   Source: %s\n", result.Profile.Name)
		}
		fmt.Println()
	}

	// Use helper functions from WebSearchResponse
	if results.HasMoreResults() {
		fmt.Println("More results are available.")
	}

	fmt.Println("\nDemonstrating search with parameters:")

	// Wait 1 second to avoid “429 Too Many Requests”
	time.Sleep(1 * time.Second)

	// Create search params with additional options
	params := &bravesearch.WebSearchParams{
		Count:      5,
		Offset:     0,
		SafeSearch: bravesearch.SafeSearchModerate,
		Freshness:  bravesearch.FreshnessWeek,
	}

	// Perform a search with parameters
	paramResults, err := client.WebSearch(ctx, query, params)
	if err != nil {
		log.Fatalf("Parametrized search failed: %v", err)
	}

	fmt.Printf("Found %d results with freshness set to past week\n\n", paramResults.GetResultCount())
	for i, result := range paramResults.GetWebResults() {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", result.URL)
		fmt.Println()
	}
}
