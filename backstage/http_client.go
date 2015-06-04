package backstage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

var BackstageClientVersion = "0.0.9"

type HTTPClient struct {
	Host   string
	client *http.Client
}

func NewHTTPClient(host string) HTTPClient {
	return HTTPClient{
		Host:   host,
		client: &http.Client{},
	}
}

type RequestArgs struct {
	AcceptableCode int
	Body           interface{}
	Path           string
	Method         string
}

func (c *HTTPClient) MakeRequest(requestArgs RequestArgs) ([]byte, error) {
	body, err := json.Marshal(requestArgs.Body)
	if err != nil {
		return []byte{}, newInvalidBodyError(err)
	}

	url, err := url.Parse(c.Host)
	if err != nil {
		return nil, newInvalidHostError(err)
	}

	url.Path = requestArgs.Path
	req, err := http.NewRequest(requestArgs.Method, url.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, newRequestError(err)
	}

	if token, err := ReadToken(); err == nil {
		req.Header.Set("Authorization", token)
	}
	req.Header.Set("BackstageClient-Version", BackstageClientVersion)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, newRequestError(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, newResponseError(err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, newUnauthorizedError(ErrLoginRequired)
	}

	if resp.StatusCode == requestArgs.AcceptableCode {
		return respBody, nil
	}

	var errorResponse ErrorResponse
	err = json.Unmarshal(respBody, &errorResponse)
	e := ErrBadResponse
	if err == nil {
		if errorResponse.Description != "" {
			e = newErrorResponse(errorResponse.Type, errorResponse.Description)
		}
	}
	return respBody, newResponseError(e)
}
