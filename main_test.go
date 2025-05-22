package main

import (
    "context"
    "testing"
    "time"
)

func TestHandleRequest(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    response, err := HandleRequest(ctx)
    if err != nil {
        t.Fatalf("Handler returned error: %v", err)
    }
    
    if response.StatusCode != 200 {
        t.Errorf("Expected status code 200, got %d", response.StatusCode)
    }
    
    t.Logf("Response body: %s", response.Body)
}
