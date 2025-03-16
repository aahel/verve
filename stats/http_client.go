package stats

import (
	"bytes"
	"net/http"
	"time"
)

// HTTPClient defines the interface for making HTTP requests
type HTTPClient interface {
	Post(url, contentType string, body []byte) (HTTPResponse, error)
}

// HTTPResponse defines the interface for HTTP responses
type HTTPResponse interface {
	StatusCode() int
	Close() error
}

// StandardHTTPClient implements HTTPClient using standard http.Client
type StandardHTTPClient struct {
	client *http.Client
}

// NewStandardHTTPClient creates a new standard HTTP client
func NewStandardHTTPClient() HTTPClient {
	return &StandardHTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Post sends a POST request
func (c *StandardHTTPClient) Post(url, contentType string, body []byte) (HTTPResponse, error) {
	resp, err := c.client.Post(url, contentType, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return &StandardHTTPResponse{resp}, nil
}

// StandardHTTPResponse implements HTTPResponse
type StandardHTTPResponse struct {
	resp *http.Response
}

// StatusCode returns the HTTP status code
func (r *StandardHTTPResponse) StatusCode() int {
	return r.resp.StatusCode
}

// Close closes the response body
func (r *StandardHTTPResponse) Close() error {
	return r.resp.Body.Close()
}
