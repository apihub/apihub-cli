package apihub_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/apihub/apihub-cli/maestro"
	"github.com/tsuru/tsuru/fs/fstest"
	. "gopkg.in/check.v1"
)

func (s *S) TestMakeRequest(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "Alice"}`))
	}))
	defer server.Close()
	httpClient.Host = server.URL

	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: nil, AcceptableCode: http.StatusOK}
	body, err := httpClient.MakeRequest(args)
	c.Assert(string(body), Equals, `{"name": "Alice"}`)
	c.Check(err, IsNil)
}

func (s *S) TestMakeRequestWithNonAcceptableCode(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "Alice"}`))
	}))
	defer server.Close()
	httpClient.Host = server.URL

	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: nil, AcceptableCode: http.StatusBadRequest}
	body, err := httpClient.MakeRequest(args)
	c.Assert(string(body), Equals, `{"name": "Alice"}`)
	e, ok := err.(apihub.ResponseError)
	c.Assert(ok, Equals, true)
	c.Assert(e.Error(), Equals, "The response was invalid or cannot be served. For more details, execute the command with `-h`.")
}

func (s *S) TestReturnsErrorWhenPayloadIsInvalid(c *C) {
	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: unsupportedPayload}
	_, err := httpClient.MakeRequest(args)
	_, ok := err.(apihub.InvalidBodyError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestReturnsErrorWhenHostIsInvalid(c *C) {
	httpClient.Host = "://invalid-host"
	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: nil}
	_, err := httpClient.MakeRequest(args)
	_, ok := err.(apihub.InvalidHostError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestReturnsErrorWhenRequestIsInvalid(c *C) {
	httpClient.Host = "invalid-host"
	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: nil}
	_, err := httpClient.MakeRequest(args)
	_, ok := err.(apihub.RequestError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestReturnsErrorWhenResponseIsInvalid(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Length", "1")
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	httpClient.Host = server.URL

	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: nil}
	_, err := httpClient.MakeRequest(args)
	_, ok := err.(apihub.ResponseError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestIncludesTokenInHeader(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "Token abcde"}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		c.Assert(auth, Equals, "Token abcde")
	}))
	defer server.Close()

	httpClient.Host = server.URL

	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: nil}
	httpClient.MakeRequest(args)
}

func (s *S) TestShouldIncludeTheClientVersionInTheHeader(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		version := req.Header.Get("ApiHubClient-Version")
		c.Assert(version, Equals, apihub.ApiHubClientVersion)
	}))
	defer server.Close()

	httpClient.Host = server.URL

	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: nil}
	httpClient.MakeRequest(args)
}

func (s *S) TestReturnsUnauthorizedError(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	httpClient.Host = server.URL

	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: nil}
	_, err := httpClient.MakeRequest(args)
	_, ok := err.(apihub.UnauthorizedError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestReturnsErrorForBadRequest(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error_description": "Something went wrong."}`))
	}))
	defer server.Close()

	httpClient.Host = server.URL

	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: nil}
	_, err := httpClient.MakeRequest(args)
	e, ok := err.(apihub.ResponseError)
	c.Assert(e.Error(), Equals, "Something went wrong.")
	c.Assert(ok, Equals, true)
}

func (s *S) TestReturnsDefaultError(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	httpClient.Host = server.URL

	args := apihub.RequestArgs{Method: "GET", Path: "/path", Body: nil}
	_, err := httpClient.MakeRequest(args)
	e, ok := err.(apihub.ResponseError)
	c.Assert(e.Error(), Equals, "The response was invalid or cannot be served. For more details, execute the command with `-h`.")
	c.Assert(ok, Equals, true)
}
