package client

import (
	"fmt"
	"net/http"
)

const apiKeysPathName string = "api_keys"
const singleKeyResponseName string = "api_key"

func (c *Client) GetApiKey(apiKeyId string) (*ApiKeyResponse, error) {
	requestOptions := RequestOptions{
		Method:       http.MethodGet,
		Path:         fmt.Sprintf("/%s/%s", apiKeysPathName, apiKeyId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[ApiKeyResponse](c, requestOptions, singleKeyResponseName)
}

func (c *Client) CreateApiKey(newApiKey ApiKeyRequest) (*ApiKeyResponse, error) {
	requestOptions := RequestOptions{
		Body:         newApiKey,
		Method:       http.MethodPost,
		Path:         fmt.Sprintf("/%s", apiKeysPathName),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[ApiKeyResponse](c, requestOptions, singleKeyResponseName)
}

func (c *Client) UpdateApiKey(apiKeyId string, updatedApiKey ApiKeyRequest) (*ApiKeyResponse, error) {
	requestOptions := RequestOptions{
		Body:         updatedApiKey,
		Method:       http.MethodPut,
		Path:         fmt.Sprintf("/%s/%s", apiKeysPathName, apiKeyId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[ApiKeyResponse](c, requestOptions, singleKeyResponseName)
}

func (c *Client) DeleteApiKey(apiKeyId string) (*ApiKeyResponse, error) {
	requestOptions := RequestOptions{
		Method:       http.MethodDelete,
		Path:         fmt.Sprintf("/%s/%s", apiKeysPathName, apiKeyId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[ApiKeyResponse](c, requestOptions, "_links") // because delete requests will not have an entity in the response
}
