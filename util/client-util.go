package util

import (
	"awsx-api/models"
	"net/http"
	"time"
)

func NewGrafanaClient() *models.GrafanaClient {
	return NewGrafanaClientWithHTTPClient(&http.Client{
		Timeout: 25 * time.Second,
	})
}

// NewGrafanaClientWithHTTPClient returns a new GrafanaClient with the given HTTP Client
func NewGrafanaClientWithHTTPClient(client *http.Client) *models.GrafanaClient {
	return &models.GrafanaClient{
		HttpClient: client,
	}
}
