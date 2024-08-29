package client

import (
	"fmt"
	"net/http"
)

const singleListenerResponseKey string = "listener"
const listenersPathName string = singleListenerResponseKey + "s"

func (c *Client) GetListener(queueId string, listenerId string) (*ListenerResponse, error) {
	requestOptions := RequestOptions{
		Method:       http.MethodGet,
		Path:         fmt.Sprintf("/queues/%s/%s/%s", queueId, listenersPathName, listenerId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[ListenerResponse](c, requestOptions, singleListenerResponseKey)
}

func (c *Client) CreateListener(queueId string, newListener ListenerRequest) (*ListenerResponse, error) {
	requestOptions := RequestOptions{
		Body:         newListener,
		Method:       http.MethodPost,
		Path:         fmt.Sprintf("/queues/%s/%s", queueId, listenersPathName),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[ListenerResponse](c, requestOptions, singleListenerResponseKey)
}

func (c *Client) UpdateListener(queueId string, listenerId string, updatedListener ListenerRequest) (*ListenerResponse, error) {
	requestOptions := RequestOptions{
		Body:         updatedListener,
		Method:       http.MethodPut,
		Path:         fmt.Sprintf("/queues/%s/%s/%s", queueId, listenersPathName, listenerId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[ListenerResponse](c, requestOptions, singleListenerResponseKey)
}

func (c *Client) DeleteListener(queueId string, listenerId string) (*ListenerResponse, error) {
	requestOptions := RequestOptions{
		Method:       http.MethodDelete,
		Path:         fmt.Sprintf("/queues/%s/%s/%s", queueId, listenersPathName, listenerId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[ListenerResponse](c, requestOptions, "_links") // because delete requests will not have an entity in the response
}
