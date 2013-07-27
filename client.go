package http

import (
	"github.com/gorilla/http/client"
	"io"
	stdurl "net/url"
	"strings"
)

// Client implements a high level HTTP client. Client methods can be called concurrently
// to as many end points as required.
// Concurrency, connection reuse, caching, and keepalive behavior is managed by the
// ConnectionManager.
type Client struct {
	dialer Dialer
}

// Do sends an HTTP request and returns an HTTP response.
func (c *Client) Do(method, url string, headers map[string][]string, body io.Reader) (client.Status, map[string][]string, io.ReadCloser, error) {
	u, err := stdurl.Parse(url)
	if err != nil {
		return client.Status{}, nil, nil, err
	}
	host := u.Host
	if !strings.Contains(host, ":") {
		host += ":80"
	}
	conn, err := c.dialer.Dial("tcp", host)
	if err != nil {
		return client.Status{}, nil, nil, err
	}
	req := client.Request{
		Method:  method,
		URI:     u.Path,
		Version: client.HTTP_1_1,
	}
	if err := conn.WriteRequest(&req); err != nil {
		return client.Status{}, nil, nil, err
	}
	resp, err := conn.ReadResponse()
	if err != nil {
		return client.Status{}, nil, nil, err
	}
	return resp.Status, nil, &readCloser{resp.Body, conn}, nil
}

type readCloser struct {
	io.Reader
	io.Closer
}

// Get sends a GET request
func (c *Client) Get(url string, headers map[string][]string) (client.Status, map[string][]string, io.ReadCloser, error) {
	return c.Do("GET", url, headers, nil)
}
