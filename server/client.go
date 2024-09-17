package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type ClientAPI struct {
	http.Client
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

	return nil
}

// SendAPIRequest Will wrap and call the sendRequest
// SendAPIRequest should be concurrent safe
func (c *ClientAPI) SendAPIRequest(url, method string, body io.Reader) error {
	ctx := context.Background()
	// Parse url

	//sendRequest()

	return nil
}
