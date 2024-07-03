package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Define a struct that matches the JSON data structure
type ApiResponse struct {
	URL string `json:"url"`
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	targetURL := r.URL.Query().Get("url")

	// Check if it's a preflight request and handle it
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK) // Send 200 OK for preflight requests
		log.Printf("Handled preflight OPTIONS request")
		return // Stop further processing
	}

	if targetURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		log.Printf("Bad request: URL parameter is missing")
		return
	}

	proxyURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Printf("Invalid URL: %s", targetURL)
		return
	}

	// Create a new request to the target URL
	proxyReq, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		log.Printf("Failed to create request for %s: %v", proxyURL.String(), err)
		return
	}

	// Copy headers from the original request to the proxy request
	proxyReq.Header = r.Header

	// Use the original request's body
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		proxyReq.Body = r.Body
	}

	// Forward the request to the target server
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("%s %s -> %s %d", r.Method, r.URL.String(), proxyURL.String(), resp.StatusCode)

	// Copy headers from the response to the client
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Write the status code
	w.WriteHeader(resp.StatusCode)

	// Add CORS headers to the response
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Copy the response body to the client
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response body: %v", err)
	}

	// if resp.StatusCode == http.StatusOK {
	// 	bodyBytes, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	bodyString := string(bodyBytes)
	// 	log.Printf(bodyString)
	// }

	duration := time.Since(start)
	log.Printf("%s %s -> %s %d (%v)", r.Method, r.URL.String(), proxyURL.String(), resp.StatusCode, duration)
}

func main() {
	http.HandleFunc("/proxy", handleProxy)
	log.Println("Starting proxy server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
