// This file is a free adaptation of Tsuru client: https://github.com/tsuru/tsuru/blob/master/cmd/client.go
package main

import (
	"bytes"
	"net/http"
	"net/url"

	httpErr "github.com/backstage/backstage/errors"
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{HTTPClient: client}
}

func (c *Client) checkTargetError(err error) error {
	urlErr, ok := err.(*url.Error)
	if !ok {
		return err
	}
	return &httpErr.HTTPError{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Failed to connect to Backstage server: " + urlErr.Err.Error(),
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if token, err := ReadToken(); err == nil {
		req.Header.Set("Authorization", token)
	}

	resp, err := c.HTTPClient.Do(req)
	err = c.checkTargetError(err)
	if err != nil {
		return nil, err
	}
	req.Close = true

	if resp.StatusCode >= 400 {
		var httpResponse = map[string]interface{}{}
		parseBody(resp.Body, &httpResponse)
		switch resp.StatusCode {
		case 401:
			err = &httpErr.HTTPError{
				StatusCode: resp.StatusCode,
				Message:    ErrLoginRequired.Error(),
			}
		default:
			err = &httpErr.HTTPError{
				StatusCode: resp.StatusCode,
				Message:    httpResponse["message"].(string),
			}
		}
		return resp, err
	}

	return resp, nil
}

func (c *Client) MakePost(path string, payload string, r interface{}) (*http.Response, error) {
	url, err := GetURL(path)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBufferString(payload)
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}

	response, err := c.Do(req)
	if err != nil {
		httpEr := err.(*httpErr.HTTPError)
		return nil, httpEr
	}

	parseBody(response.Body, &r)
	return response, nil
}

func (c *Client) MakeDelete(path string, payload string, r interface{}) (*http.Response, error) {
	url, err := GetURL(path)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBufferString(payload)
	req, err := http.NewRequest("DELETE", url, b)
	if err != nil {
		return nil, err
	}

	response, err := c.Do(req)
	if err != nil {
		httpEr := err.(*httpErr.HTTPError)
		return nil, httpEr
	}

	parseBody(response.Body, &r)
	return response, nil
}
