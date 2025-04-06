package events

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// DefaultHttpClient implements the HttpClient interface
type DefaultHttpClient struct {
	client *http.Client
}

// NewHttpClient creates a new HTTP client
func NewHttpClient(timeout time.Duration, insecure bool) *DefaultHttpClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	return &DefaultHttpClient{
		client: &http.Client{
			Timeout:   timeout,
			Transport: tr,
		},
	}
}

// Get performs an HTTP GET request
func (c *DefaultHttpClient) Get(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status %d: %s",
			resp.StatusCode, body)
	}

	return body, nil
}

// Post performs an HTTP POST request
func (c *DefaultHttpClient) Post(url string, contentType string, body []byte) ([]byte, error) {
	resp, err := c.client.Post(url, contentType, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("HTTP POST request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status %d: %s",
			resp.StatusCode, respBody)
	}

	return respBody, nil
}

// Do performs a custom HTTP request
func (c *DefaultHttpClient) Do(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status %d: %s",
			resp.StatusCode, body)
	}

	return body, nil
}
