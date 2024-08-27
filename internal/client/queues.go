package client

import (
	"fmt"
	"net/http"
)

const queuesResourceName string = "queues"
const singleQueueResponseKey string = "queue"

func (c *Client) GetQueue(queueId string) (*Queue, error) {
	requestOptions := RequestOptions{
		Method:       http.MethodGet,
		Path:         fmt.Sprintf("/%s/%s", queuesResourceName, queueId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[Queue](c, requestOptions, singleQueueResponseKey)
}

func (c *Client) CreateQueue(newQueue Queue) (*Queue, error) {
	requestOptions := RequestOptions{
		Body:         newQueue,
		Method:       http.MethodPost,
		Path:         fmt.Sprintf("/%s", queuesResourceName),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[Queue](c, requestOptions, singleQueueResponseKey)
}

func (c *Client) UpdateQueue(queueId string, updatedQueue Queue) (*Queue, error) {
	requestOptions := RequestOptions{
		Body:         updatedQueue,
		Method:       http.MethodPut,
		Path:         fmt.Sprintf("/%s/%s", queuesResourceName, queueId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[Queue](c, requestOptions, singleQueueResponseKey)
}

func (c *Client) DeleteQueue(queueId string) (*Queue, error) {
	requestOptions := RequestOptions{
		Method:       http.MethodDelete,
		Path:         fmt.Sprintf("/%s/%s", queuesResourceName, queueId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[Queue](c, requestOptions, "_links") // because delete requests will not have an entity in the response
}
