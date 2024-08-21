package client

import "net/http"

type Client struct {
	ApiEndpoint string
	ApiKey      string
	HttpClient  *http.Client
}

type RequestOptions struct {
	Method       string
	Path         string
	ExpectStatus int
	Body         any
}

type Queue struct {
	Id    string `json:"id,omitempty"`
	Alias string `json:"alias"`
}

type DomainRequest struct {
	Alias    string `json:"alias"`
	Hostname string `json:"hostname,omitempty"`
	Port     int    `json:"port,omitempty"`
}

type DomainResponse struct {
	Id           string              `json:"id"`
	Alias        string              `json:"alias"`
	Hostname     string              `json:"hostname"`
	Port         int32               `json:"port"`
	Challenge    *DomainChallenge    `json:"challenge"`
	Verification *DomainVerification `json:"verification"`
}

type DomainChallenge struct {
	Https HttpsDomainChallenge `json:"https"`
}

type HttpsDomainChallenge struct {
	FileContent string `json:"file_content"`
	FileName    string `json:"file_name"`
	ContextPath string `json:"context_path"`
	CreatedAt   int64  `json:"created_at"`
	ExpiresAt   int64  `json:"expires_at"`
}

type DomainVerification struct {
	Verified   bool  `json:"verified,omitempty"`
	VerifiedAt int64 `json:"verified_at,omitempty"`
}

type ApiResponse[T any] struct {
	Response T `json:"-"` // Use a custom JSON key
}
