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

type ApiResource struct {
	Id    string `json:"id,omitempty"`
	Alias string `json:"alias"`
}

type Queue struct {
	*ApiResource
}

type QueueReponse struct {
	Queue any `json:"queue"`
}
