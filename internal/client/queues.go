package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const queues string = "queues"

func sendAndUnmarshal[T any](c *Client, requestOptions RequestOptions) (*T, error) {
	body, err := c.executeRequest(requestOptions)
	if err != nil {
		return nil, err
	}

	response := new(T)
	err = json.Unmarshal(body, &QueueResponse{response})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetQueue(queueId string) (*Queue, error) {
	requestOptions := RequestOptions{
		Method:       http.MethodGet,
		Path:         fmt.Sprintf("/%s/%s", queues, queueId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndUnmarshal[Queue](c, requestOptions)
}

func (c *Client) CreateQueue(newQueue Queue) (*Queue, error) {
	requestOptions := RequestOptions{
		Body:         newQueue,
		Method:       http.MethodPost,
		Path:         fmt.Sprintf("/%s", queues),
		ExpectStatus: http.StatusOK,
	}

	return sendAndUnmarshal[Queue](c, requestOptions)
}

func (c *Client) UpdateQueue(queueId string, updatedQueue Queue) (*Queue, error) {
	requestOptions := RequestOptions{
		Body:         updatedQueue,
		Method:       http.MethodPut,
		Path:         fmt.Sprintf("/%s/%s", queues, queueId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndUnmarshal[Queue](c, requestOptions)
}

func (c *Client) DeleteQueue(queueId string) (*Queue, error) {
	requestOptions := RequestOptions{
		Method:       http.MethodDelete,
		Path:         fmt.Sprintf("/%s/%s", queues, queueId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndUnmarshal[Queue](c, requestOptions)
}
