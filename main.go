package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

// Post represents the structure of data from JSONPlaceholder API
type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Response is the structure we'll return from our Lambda function
type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// fetchPosts retrieves posts from the JSONPlaceholder API
func fetchPosts(ctx context.Context) ([]Post, error) {
	// Get API URL from environment variable or use default
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "https://jsonplaceholder.typicode.com/posts"
	}

	log.Printf("Fetching data from: %s", apiURL)

	// Create an HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new request with context for cancellation support
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned non-OK status: %d", resp.StatusCode)
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal JSON response
	var posts []Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return posts, nil
}

func fetchPostsWithRetry(ctx context.Context, maxRetries int) ([]Post, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		posts, err := fetchPosts(ctx)
		if err == nil {
			return posts, nil
		}

		lastErr = err
		log.Printf("Retry %d/%d failed: %v", i+1, maxRetries, err)

		// Exponential backoff
		sleepTime := time.Duration((1<<i)*100) * time.Millisecond
		select {
		case <-time.After(sleepTime):
			continue
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled during retry: %v", ctx.Err())
		}
	}

	return nil, fmt.Errorf("all %d retries failed, last error: %v", maxRetries, lastErr)
}

// HandleRequest is our Lambda function handler
func HandleRequest(ctx context.Context) (Response, error) {
	log.Println("Lambda execution started")

	// Get max retries from environment variable or use default
	maxRetries := 3
	if maxRetriesStr := os.Getenv("MAX_RETRIES"); maxRetriesStr != "" {
		if val, err := strconv.Atoi(maxRetriesStr); err == nil {
			maxRetries = val
		}
	}

	// Fetch posts from API
	posts, err := fetchPostsWithRetry(ctx, maxRetries) // Try up to 3 times

	if err != nil {
		log.Printf("Failed to fetch posts: %v", err)
		return Response{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf(`{"error": "Failed to fetch posts: %v"}`, err),
		}, nil
	}

	// Log the total number of items
	totalItems := len(posts)
	log.Printf("Total items fetched: %d", totalItems)

	// Check if we have at least one item
	if totalItems == 0 {
		log.Println("No items returned from the API")
		return Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"message": "No items found", "total": 0}`,
		}, nil
	}

	// Log the title of the first item
	firstItemTitle := posts[0].Title
	log.Printf("Title of first item: %s", firstItemTitle)

	// Create response body
	responseData := map[string]interface{}{
		"total":          totalItems,
		"firstItemTitle": firstItemTitle,
	}

	responseBody, err := json.Marshal(responseData)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return Response{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "Failed to generate response"}`,
		}, nil
	}

	log.Println("Lambda execution completed successfully")
	return Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseBody),
	}, nil
}

func main() {
	// Start the Lambda handler
	lambda.Start(HandleRequest)
}
