package workwave

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	qt "github.com/frankban/quicktest"
)

var (
	mux    *http.ServeMux
	ctx    = context.TODO()
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

func teardown() {
	server.Close()
}

func TestNew(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		c := qt.New(t)
		client, err := New("api-key")
		c.Assert(err, qt.IsNil)

		clientURL, _ := url.Parse(apiBaseURL)
		c.Assert(client.baseURL, qt.DeepEquals, clientURL)
		c.Assert(client.apiKey, qt.Equals, "api-key")
		c.Assert(client.client, qt.Not(qt.IsNil))
	})
}

func TestNewRequest(t *testing.T) {
	t.Run("valid, no body", func(t *testing.T) {
		c := qt.New(t)
		client, err := New("api-key")
		c.Assert(err, qt.IsNil)

		req, err := client.NewRequest(ctx, http.MethodGet, "/path/", nil)
		c.Assert(err, qt.IsNil)
		// Method
		c.Assert(req.Method, qt.Equals, http.MethodGet)
		// Headers
		c.Assert(req.Header.Get("X-WorkWave-Key"), qt.Equals, "api-key")
		c.Assert(req.Header.Get("User-Agent"), qt.Equals, agentString)
		c.Assert(req.Header.Get("Content-Type"), qt.Equals, contentType)
		c.Assert(req.Header.Get("Accept"), qt.Equals, contentType)
		// URL
		URL, _ := url.Parse(apiBaseURL + "/path/")
		c.Assert(req.URL, qt.DeepEquals, URL)
	})

	t.Run("valid, with body", func(t *testing.T) {
		c := qt.New(t)
		client, err := New("api-key")
		c.Assert(err, qt.IsNil)

		b := `{"json": "data"}`
		req, err := client.NewRequest(ctx, http.MethodGet, "/path/", b)
		c.Assert(err, qt.IsNil)
		// Method
		c.Assert(req.Method, qt.Equals, http.MethodGet)
		// Headers
		c.Assert(req.Header.Get("X-WorkWave-Key"), qt.Equals, "api-key")
		c.Assert(req.Header.Get("User-Agent"), qt.Equals, agentString)
		c.Assert(req.Header.Get("Content-Type"), qt.Equals, contentType)
		c.Assert(req.Header.Get("Accept"), qt.Equals, contentType)
		// URL
		URL, _ := url.Parse(apiBaseURL + "/path/")
		c.Assert(req.URL, qt.DeepEquals, URL)
		// Body
		rBody, _ := ioutil.ReadAll(req.Body)
		c.Assert(rBody, qt.JSONEquals, string(b))
	})

	t.Run("invalid method", func(t *testing.T) {
		c := qt.New(t)
		client, err := New("api-key")
		c.Assert(err, qt.IsNil)

		_, err = client.NewRequest(ctx, "{", "/path/", nil)
		c.Assert(err, qt.ErrorMatches, ".*invalid method.*")
	})

	t.Run("invalid path", func(t *testing.T) {
		c := qt.New(t)
		client, err := New("api-key")
		c.Assert(err, qt.IsNil)

		_, err = client.NewRequest(ctx, http.MethodGet, "%zz", nil)
		c.Assert(err, qt.ErrorMatches, ".*invalid URL escape.*")
	})

	t.Run("invalid body", func(t *testing.T) {
		c := qt.New(t)
		client, err := New("api-key")
		c.Assert(err, qt.IsNil)

		b := make(chan int)
		_, err = client.NewRequest(ctx, http.MethodGet, "/path/", b)
		c.Assert(err, qt.ErrorMatches, ".*unsupported type.*")
	})
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()
	c := qt.New(t)

	mux.HandleFunc("/good", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"json":"data"}`)
	})
	mux.HandleFunc("/bad", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	})

	client, _ := New("api-key")
	client.baseURL, _ = url.Parse(server.URL)

	t.Run("valid", func(t *testing.T) {
		req, _ := client.NewRequest(ctx, http.MethodGet, "/good", nil)
		_, err := client.Do(ctx, req, nil)
		c.Assert(err, qt.IsNil)
	})

	t.Run("valid with body return", func(t *testing.T) {
		type stuff struct {
			JSON string
		}
		s := new(stuff)
		req, _ := client.NewRequest(ctx, http.MethodGet, "/good", s)
		_, err := client.Do(ctx, req, s)
		c.Assert(err, qt.IsNil)
		c.Assert(s, qt.DeepEquals, &stuff{JSON: "data"})
	})

	t.Run("invalid body return type", func(t *testing.T) {

		req, _ := client.NewRequest(ctx, http.MethodGet, "/good", 1)
		_, err := client.Do(ctx, req, 1)
		c.Assert(err, qt.ErrorMatches, ".*json: Unmarshal.*")
	})

	t.Run("bad response", func(t *testing.T) {
		req, _ := client.NewRequest(ctx, http.MethodGet, "/bad", nil)
		_, err := client.Do(ctx, req, nil)
		c.Assert(err, qt.ErrorMatches, "HTTP.*error")
	})

}

func TestCheckResponse(t *testing.T) {
	for _, tt := range []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "200 range",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "< 200",
			statusCode: http.StatusContinue,
			wantErr:    true,
		},
		{
			name:       "> 300",
			statusCode: http.StatusBadGateway,
			wantErr:    true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			resp := &http.Response{StatusCode: tt.statusCode}
			err := checkResponse(resp)
			if tt.wantErr {
				c.Assert(err, qt.ErrorMatches, "HTTP.*error")
			} else {
				c.Assert(err, qt.IsNil)
			}
		})
	}
}
