package workwave

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

const (
	agentString = "patch-gardens/go-workwave"
	contentType = "application/json"
	apiBaseURL  = "https://wwrm.workwave.com"
)

// Client is a structure that provides access to the WorkWave API.
// https://wwrm.workwave.com/api/
type Client struct {
	client  *http.Client
	baseURL *url.URL
	apiKey  string
}

// New creates a new WorkWave API client with the given API key for authentication.
func New(apiKey string) (*Client, error) {
	baseURL, err := url.Parse(apiBaseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		client:  newHTTPClient(),
		baseURL: baseURL,
		apiKey:  apiKey,
	}
	return c, nil
}

// newHTTPClient creates an http.Client with timeouts which will be used to make
// requests to the Clubhouse API.
func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 2 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
}

// NewRequest prepares an API HTTP request using the given method and path.
// If body is given, it will be converted to JSON and included in the request.
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Add HTTP headers
	req.Header.Add("X-WorkWave-Key", c.apiKey)
	req.Header.Add("User-Agent", agentString)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", contentType)
	return req, nil
}

// Do submits an HTTP request with the Client's HTTP client.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := res.Body.Close(); err == nil {
			err = rerr
		}
	}()

	err = checkResponse(res)
	if err != nil {
		return res, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, res.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(res.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}
	return res, err
}

func checkResponse(res *http.Response) error {
	sc := res.StatusCode
	if sc >= 200 && sc <= 200 {
		return nil
	}
	return fmt.Errorf("HTTP %d error", sc)
}
