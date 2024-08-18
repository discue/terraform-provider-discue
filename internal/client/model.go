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

type QueueResponse struct {
	Queue any `json:"queue"`
}
