package server

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type ClientContext struct {
	*http.Client
	ClientCount  int
	RequestCount int // basic right now with no rate limiting
}

// NewClientContext initializes ClientContext with a default HTTP client
func NewClientContext() *ClientContext {
	return &ClientContext{
		&http.Client{Timeout: 5 * time.Second}, // Example timeout
		0,
		0,
	}
}

type User struct {
	Name           string
	Token          string
	RequestHistory *RequestHistory
	mu             sync.Mutex
}

type RequestHistory struct {
	Request   *Request
	TimeStamp time.Time
}

type ClientAPI struct {
	Client      *ClientContext
	User        *User
	RequestInfo *Request
	mu          sync.Mutex // Protect RequestInfo from concurrent access
}

// Request should be used as value and not changed during a client request
// The struct will hold request details for the client request and add these to a pointed value for a User and their
// RequestHistory
type Request struct {
	URL     string
	Method  string
	Body    ReqBody
	Headers map[string]string
}

type ReqBody struct {
	Body io.Reader
}

// SendAPIRequest will wrap and call the sendRequest
// SendAPIRequest should be concurrent safe
func (c *ClientAPI) SendAPIRequest(url, method string, body io.Reader) {
	// Locking to safely update RequestInfo
	c.mu.Lock()
	defer c.mu.Unlock()

	// Initialize RequestInfo
	c.RequestInfo = &Request{
		URL:    url,
		Method: method,
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}

	// Only set body for methods that require it
	if method == http.MethodPost || method == http.MethodPut {
		c.RequestInfo.Body = ReqBody{Body: body} // Set body if needed.
	}

	// Send the request in a goroutine
	go func() {
		fmt.Printf("Attempting to send request: %+v\n", c.RequestInfo) // Debug logging
		err := sendRequest(c)                                          // Call sendRequest without context
		if err != nil {
			fmt.Printf("error sending request: %v\n", err)
		}
	}()
}

// sendRequest sends the HTTP request and logs the response
func sendRequest(c *ClientAPI) error {
	fmt.Printf("Sending request to URL: %s with method: %s\n", c.RequestInfo.URL, c.RequestInfo.Method)

	if c.RequestInfo == nil {
		return fmt.Errorf("RequestInfo is nil")
	}
	fmt.Printf("RequestInfo: %+v\n", c.RequestInfo)

	req, err := http.NewRequest(c.RequestInfo.Method, c.RequestInfo.URL, c.RequestInfo.Body.Body)
	if err != nil {
		return fmt.Errorf("error forming request: %v", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "go-proxpost-Runtime")
	for key, value := range c.RequestInfo.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read and log response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Host: %s\n", req.UserAgent())
	fmt.Printf("Response body: %s\n", string(body))

	return nil
}
