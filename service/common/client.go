package common

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type RestClient struct {
	baseUrl       *url.URL
	client        *http.Client
	Header        *http.Header
	contentConfig ContentConfig
}

func (rc *RestClient) Get() *Request {
	return rc.Verb("GET")
}

func (rc *RestClient) POST() *Request {
	return rc.Verb("POST")
}

func (rc *RestClient) PUT() *Request {
	return rc.Verb("PUT")
}

func (rc *RestClient) DELETE() *Request {
	return rc.Verb("DELETE")
}

func (rc *RestClient) HEAD() *Request {
	return rc.Verb("HEAD")
}

func (rc *RestClient) Verb(verb string) *Request {
	return NewRequest(rc.client, verb, rc.baseUrl, rc.contentConfig, rc.client.Timeout)
}

func NewRestClient(baseUrl *url.URL, timeout int64) *RestClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &RestClient{
		baseUrl: baseUrl,
		client: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(timeout) * time.Second,
		},
	}
}
