package workwave

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestCallbackGet(t *testing.T) {
	setup()
	defer teardown()
	c := qt.New(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"url": "https://my.server.com/callback"}`)
	})

	client, _ := New("api-key")
	client.baseURL, _ = url.Parse(server.URL)

	cb, err := client.Callback.Get(ctx, Callback{
		URL: "https://my.server.com/callback",
	})
	c.Assert(err, qt.IsNil)
	c.Assert(cb.URL, qt.Equals, "https://my.server.com/callback")
}

func TestCallbackSet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var rBody Callback
		if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
			t.Fatal("failed to parse request body")
		}
		switch {
		case rBody.Test: // setting Test returns an error to check for it
			fmt.Fprintf(w, `{"errorCode": 2000, "errorMessage": "Server at URL [https://my.server.com/callback] failed to respond to the test message."}`)
			return
		default:
			fmt.Fprintf(w, `{"url":"https://my.server.com/new-callback", "previousUrl":"https://my.server.com/callback"}`)
			return
		}
	})

	client, _ := New("api-key")
	client.baseURL, _ = url.Parse(server.URL)

	t.Run("simple", func(t *testing.T) {
		c := qt.New(t)
		cb, err := client.Callback.Set(ctx, Callback{
			URL: "https://my.server.com/new-callback",
		})
		c.Assert(err, qt.IsNil)
		c.Assert(cb.URL, qt.Equals, "https://my.server.com/new-callback")
		c.Assert(cb.PreviousURL, qt.Equals, "https://my.server.com/callback")
	})

	t.Run("with test and error response", func(t *testing.T) {
		c := qt.New(t)
		cb, err := client.Callback.Set(ctx, Callback{
			URL:  "https://my.server.com/new-callback",
			Test: true,
		})
		c.Assert(err, qt.ErrorMatches, "failed to set callback.*")
		c.Assert(cb, qt.DeepEquals, Callback{
			ErrorCode:    2000,
			ErrorMessage: "Server at URL [https://my.server.com/callback] failed to respond to the test message.",
		})
	})
}

func TestCallbackDelete(t *testing.T) {
	setup()
	defer teardown()
	c := qt.New(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"previousUrl": "https://my.server.com/callback"}`)
	})

	client, _ := New("api-key")
	client.baseURL, _ = url.Parse(server.URL)

	cb, err := client.Callback.Delete(ctx, Callback{
		URL: "https://my.server.com/callback",
	})
	c.Assert(err, qt.IsNil)
	c.Assert(cb.PreviousURL, qt.Equals, "https://my.server.com/callback")
}
