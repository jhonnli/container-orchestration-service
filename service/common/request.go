package common

import (
	"bytes"
	"fmt"
	"github.com/jhonnli/logs"
	"github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Request struct {
	client HTTPClient
	verb   string

	baseURL *url.URL
	timeout time.Duration
	content ContentConfig
	path    string
	params  url.Values
	headers http.Header

	err  error
	body io.Reader
}

func (req *Request) RequestURI(path string) *Request {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	req.path = path
	return req
}

func (req *Request) URL() *url.URL {

	p := req.baseURL.Path
	if req.path != "" {
		p += req.path
	}
	finalURL := &url.URL{}
	if req.baseURL != nil {
		*finalURL = *req.baseURL
	}
	finalURL.Path = p
	query := url.Values{}
	for key, values := range req.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	if req.timeout != 0 {
		query.Set("timeout", req.timeout.String())
	}
	finalURL.RawQuery = query.Encode()
	return finalURL
}

func (req *Request) Body(obj interface{}) *Request {
	if req.err != nil {
		return req
	}
	switch t := obj.(type) {
	case string:
		data, err := ioutil.ReadFile(t)
		if err != nil {
			req.err = err
			return req
		}
		req.body = bytes.NewBuffer(data)
	case []byte:
		req.body = bytes.NewBuffer(t)
	case io.Reader:
		req.body = t
	default:
		req.err = fmt.Errorf("unknow type used for body: %+v", obj)
	}
	return req
}

func (r *Request) Param(paramName, s string) *Request {
	if r.err != nil {
		return r
	}
	return r.setParam(paramName, s)
}

func (r *Request) setParam(paramName, value string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[paramName] = append(r.params[paramName], value)
	return r
}

func (req *Request) Do() Result {
	var result Result
	if req.client == nil {
		req.client = http.DefaultClient
	}

	url := req.URL().String()
	request, err := http.NewRequest(req.verb, url, req.body)
	if err != nil {
		return Result{err: err}
	}
	request.Header = req.headers
	resp, err := req.client.Do(request)
	if err != nil {
		logs.Error(err.Error())
		if err, ok := err.(net.Error); ok && err.Timeout() {
			req.err = err
			return Result{
				statusCode: http.StatusRequestTimeout,
			}
		}
	}
	data, err := ioutil.ReadAll(resp.Body)
	log.Println("Value: ", string(data))
	resp.Body.Close()
	result = Result{
		body:        data,
		statusCode:  resp.StatusCode,
		contentType: resp.Header.Get("Content-Type"),
	}
	return result
}

func NewRequest(client HTTPClient, verb string, baseURL *url.URL, content ContentConfig, timeout time.Duration) *Request {

	return &Request{
		client:  client,
		verb:    verb,
		baseURL: baseURL,
		timeout: timeout,
		content: content,
	}
}

type Result struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
}

func (r Result) StatusCode(statusCode *int) Result {
	*statusCode = r.statusCode
	return r
}

func (r Result) Into(obj interface{}) error {
	log.Println(string(r.body))
	err := json.Unmarshal(r.body, obj)
	return err
}

type ContentConfig struct {
	AcceptContentTypes string
	ContentType        string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Values map[string][]string

func (v Values) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (v Values) Set(key, value string) {
	v[key] = []string{value}
}

func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}

func (v Values) Del(key string) {
	delete(v, key)
}
