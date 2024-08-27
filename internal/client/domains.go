package client

import (
	"fmt"
	"net/http"
)

const domainsPathName string = "domains"
const singleDomainResponseKey string = "domain"

func (c *Client) GetDomain(domainId string) (*DomainResponse, error) {
	requestOptions := RequestOptions{
		Method:       http.MethodGet,
		Path:         fmt.Sprintf("/%s/%s", domainsPathName, domainId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[DomainResponse](c, requestOptions, singleDomainResponseKey)
}

func (c *Client) CreateDomain(newDomain DomainRequest) (*DomainResponse, error) {
	requestOptions := RequestOptions{
		Body:         newDomain,
		Method:       http.MethodPost,
		Path:         fmt.Sprintf("/%s", domainsPathName),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[DomainResponse](c, requestOptions, singleDomainResponseKey)
}

func (c *Client) UpdateDomain(domainId string, updatedDomain DomainRequest) (*DomainResponse, error) {
	requestOptions := RequestOptions{
		Body:         updatedDomain,
		Method:       http.MethodPut,
		Path:         fmt.Sprintf("/%s/%s", domainsPathName, domainId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[DomainResponse](c, requestOptions, singleDomainResponseKey)
}

func (c *Client) DeleteDomain(domainId string) (*DomainResponse, error) {
	requestOptions := RequestOptions{
		Method:       http.MethodDelete,
		Path:         fmt.Sprintf("/%s/%s", domainsPathName, domainId),
		ExpectStatus: http.StatusOK,
	}

	return sendAndReceive[DomainResponse](c, requestOptions, "_links") // because delete requests will not have an entity in the response
}
