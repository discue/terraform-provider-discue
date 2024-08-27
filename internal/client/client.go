package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func NewClient(apiEndpoint string, apiKey *string) (*Client, error) {
	c := Client{
		HttpClient:  &http.Client{Timeout: 10 * time.Second},
		ApiEndpoint: apiEndpoint,
		ApiKey:      *apiKey,
	}

	return &c, nil
}

func sendAndReceive[T any](c *Client, requestOptions RequestOptions, jsonKey string) (*T, error) {
	body, err := c.executeRequest(requestOptions)
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %w", err)
	}

	var rawMap map[string]json.RawMessage
	err = json.Unmarshal(body, &rawMap)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	if rawResponse, ok := rawMap[jsonKey]; ok {
		// Create a pointer to the response type and unmarshal into it
		responsePtr := new(T)
		err = json.Unmarshal(rawResponse, responsePtr)

		if err != nil {
			return nil, fmt.Errorf("could not unmarshal response: %w %s %s", err, rawMap, rawResponse)
		}
		return responsePtr, nil
	}

	return nil, fmt.Errorf("could not unmarshal resource from response: %s", rawMap)
}

func (c *Client) executeRequest(requestOptions RequestOptions) ([]byte, error) {
	var br io.Reader = nil
	if requestOptions.Body != nil {
		bytes, err := json.Marshal(requestOptions.Body)
		if err != nil {
			return nil, err
		}
		br = strings.NewReader(string(bytes))
	}

	req, err := http.NewRequest(requestOptions.Method, fmt.Sprintf("%s%s", c.ApiEndpoint, requestOptions.Path), br)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", c.ApiKey)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != requestOptions.ExpectStatus {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
