package models

import "time"

// Project represents a monitored project
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	RepoURL     string    `json:"repo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Trace represents a log/trace entry
type Trace struct {
	ID          string            `json:"id"`
	ProjectID   string            `json:"project_id"`
	Level       string            `json:"level"`
	Message     string            `json:"message"`
	Timestamp   time.Time         `json:"timestamp"`
	Source      string            `json:"source"`
	StackTrace  string            `json:"stack_trace"`
	Context     map[string]string `json:"context"`
	ServiceName string            `json:"service_name"`
	Environment string            `json:"environment"`
	CreatedAt   time.Time         `json:"created_at"`
}

// QueryTracesResponse represents the response from querying traces
type QueryTracesResponse struct {
	Traces []*Trace `json:"traces"`
	Total  int      `json:"total"`
}

// ListProjectsResponse represents the response from listing projects
type ListProjectsResponse struct {
	Projects []*Project `json:"projects"`
	Total    int        `json:"total"`
}

// Metrics represents service metrics
type Metrics struct {
	TotalTraces int     `json:"total_traces"`
	ErrorRate   float64 `json:"error_rate"`
	LastUpdate  string  `json:"last_update"`
}
