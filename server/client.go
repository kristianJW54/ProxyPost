package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type ClientContext struct {
	*http.Client
}

// NewClientContext Initialize ClientContext with a default HTTP client
func NewClientContext() *ClientContext {
	return &ClientContext{
		&http.Client{Timeout: 5 * time.Second}, // Example timeout
	}
}

// TODO move User to it's own module and handle

type User struct {
	//User with token and permissions
	//Session management
	//Individual rate limiting?
	RequestHistory []string
	mu             sync.Mutex
}

// ClientAPI will populate with the user and their request via middleware prior to sending the client request
type ClientAPI struct {
	client      *ClientContext
	User        *User
	RequestInfo Request
}

type Request struct {
	URL     string
	Method  string
	Body    io.Reader
	Headers map[string]string
}

//Request data structure with everything in

// un-exported functions above to provide logic for client requests to use in the below Compiled Methods

func sendRequest(ctx context.Context, c *ClientAPI) error {

	// Setting up the request
	req, err := http.NewRequestWithContext(ctx, c.RequestInfo.Method, c.RequestInfo.URL, c.RequestInfo.Body)
	if err != nil {
		return fmt.Errorf("error forming request: %v", err)
	}

	if len(c.RequestInfo.Headers) == 0 {
		//set basic headers for request
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	req.Header.Set("User-Agent", "go-proxpost-Runtime")
	for key, value := range c.RequestInfo.Headers {
		req.Header.Set(key, value)
	}

	// Continue with request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	fmt.Printf("Response: %s\n", string(body))

	// Mutex lock and add request information to user data and history
	// Also apply rate limiting logic here...?

	return nil
}

//TODO Consider SendAPIRequest taking a context and channel to collect response...

// SendAPIRequest Will wrap and call the sendRequest
// SendAPIRequest should be concurrent safe
func (c *ClientAPI) SendAPIRequest(url, method string, body io.Reader) {

	// Will receive the dependencies of the user from within the middleware function
	// to be handled here before sending request

	//Dependencies for the sendRequest internal method to run
	//Maybe user and session...?

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update request information.
	c.RequestInfo = Request{
		URL:    url,
		Method: method,
		Body:   body, // Set body if needed.
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	//sendRequest()
	go func() {
		err := sendRequest(ctx, c)
		if err != nil {
			fmt.Errorf("error sending request: %v", err)
		}
	}()
}
