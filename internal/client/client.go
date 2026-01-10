package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tapiaw38/tracehub-cli/pkg/models"
)

// Client is the HTTP client for TraceHub server
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// New creates a new TraceHub client
func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ListProjects lists all projects
func (c *Client) ListProjects(limit, offset int) (*models.ListProjectsResponse, error) {
	u := fmt.Sprintf("%s/api/v1/projects?limit=%d&offset=%d", c.baseURL, limit, offset)

	resp, err := c.doRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.ListProjectsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetProject gets a project by ID
func (c *Client) GetProject(projectID string) (*models.Project, error) {
	u := fmt.Sprintf("%s/api/v1/projects/%s", c.baseURL, projectID)

	resp, err := c.doRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var project models.Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &project, nil
}

// QueryTraces queries traces with filters
func (c *Client) QueryTraces(projectID string, filters map[string]string) (*models.QueryTracesResponse, error) {
	params := url.Values{}
	params.Add("project_id", projectID)

	for key, value := range filters {
		if value != "" {
			params.Add(key, value)
		}
	}

	u := fmt.Sprintf("%s/api/v1/traces?%s", c.baseURL, params.Encode())

	resp, err := c.doRequestWithAuth("GET", u, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.QueryTracesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Health checks server health
func (c *Client) Health() (bool, error) {
	u := fmt.Sprintf("%s/api/v1/health", c.baseURL)

	resp, err := c.doRequest("GET", u, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// doRequest performs an HTTP request without authentication
func (c *Client) doRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// doRequestWithAuth performs an HTTP request with API key authentication
func (c *Client) doRequestWithAuth(method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
