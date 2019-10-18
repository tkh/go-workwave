package workwave

import (
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestRoutesListCurrent(t *testing.T) {
	setup()
	defer teardown()
	c := qt.New(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("testdata", "routes-list-current.json"))
	})

	client, _ := New("api-key")
	client.baseURL, _ = url.Parse(server.URL)

	o, err := client.Routes.ListCurrent(ctx, RoutesListCurrentInput{
		TerritoryID: "territory",
		Date:        "20191019",
		Vehicle:     "vehicle",
	})
	c.Assert(err, qt.IsNil)
	c.Assert(len(o), qt.Equals, 2)
}

func TestRoutesListApproved(t *testing.T) {
	setup()
	defer teardown()
	c := qt.New(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("testdata", "routes-list-approved.json"))
	})

	client, _ := New("api-key")
	client.baseURL, _ = url.Parse(server.URL)

	o, err := client.Routes.ListApproved(ctx, RoutesListApprovedInput{
		TerritoryID: "territory",
		Date:        "20191019",
	})
	c.Assert(err, qt.IsNil)
	c.Assert(len(o), qt.Equals, 2)
}
