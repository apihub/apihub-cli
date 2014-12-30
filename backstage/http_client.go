package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(client *http.Client) *HTTPClient {
	return &HTTPClient{client: client}
}

func (c *HTTPClient) checkTargetError(err error) error {
	urlErr, ok := err.(*url.Error)
	if !ok {
		return err
	}

	return &HTTPError{ErrorDescription: "Failed to connect to Backstage server: " + urlErr.Err.Error()}
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("BackstageClient-Version", BackstageClientVersion)
	if token, err := ReadToken(); err == nil {
		req.Header.Set("Authorization", token)
	}
	resp, err := c.client.Do(req)
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
			err = &HTTPError{ErrorDescription: ErrLoginRequired.Error()}
		default:
			msg, ok := httpResponse["error_description"].(string)
			if !ok {
				msg = ErrFailedConnectingServer.Error()
			}

			err = &HTTPError{ErrorDescription: msg}
		}
	}

	return resp, err
}

func convertInString(p interface{}) (string, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func (c *HTTPClient) MakePost(path string, p interface{}, r interface{}) (*http.Response, error) {
	body, err := convertInString(p)
	if err != nil {
		return nil, err
	}

	url, err := GetURL(path)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBufferString(body)
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}

	response, err := c.Do(req)
	if err != nil {
		httpEr := err.(*HTTPError)
		return nil, httpEr
	}
	parseBody(response.Body, &r)
	return response, nil
}

func (c *HTTPClient) MakeDelete(path string, p interface{}, r interface{}) (*http.Response, error) {
	body, err := convertInString(p)
	if err != nil {
		return nil, err
	}

	url, err := GetURL(path)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBufferString(body)
	req, err := http.NewRequest("DELETE", url, b)
	if err != nil {
		return nil, err
	}

	response, err := c.Do(req)
	if err != nil {
		httpEr := err.(*HTTPError)
		return nil, httpEr
	}
	parseBody(response.Body, &r)
	return response, nil
}

func (c *HTTPClient) MakeGet(path string, r interface{}) (*http.Response, error) {
	url, err := GetURL(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.Do(req)
	if err != nil {
		httpEr := err.(*HTTPError)
		return nil, httpEr
	}
	parseBody(response.Body, &r)
	return response, nil
}

func GetURL(path string) (string, error) {
	t, err := LoadTargets()
	if err != nil {
		return "", err
	}
	current := t.Options[t.Current]
	if current == "" {
		return "", ErrEndpointNotFound
	}

	return strings.TrimRight(current, "/") + path, nil
}
