package workwave

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestOrdersList(t *testing.T) {
	setup()
	defer teardown()
	c := qt.New(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("testdata", "orders-list.json"))
	})

	client, _ := New("api-key")
	client.baseURL, _ = url.Parse(server.URL)

	o, err := client.Orders.List(ctx, OrdersListInput{
		TerritoryID: "territory",
		Include:     "assigned",
		EligibleOn:  "20191019",
		AssignedOn:  "20191018",
	})
	c.Assert(err, qt.IsNil)
	c.Assert(len(o), qt.Equals, 7)
}

func TestOrdersGet(t *testing.T) {
	setup()
	defer teardown()
	c := qt.New(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("testdata", "orders-get.json"))
	})

	client, _ := New("api-key")
	client.baseURL, _ = url.Parse(server.URL)

	o, err := client.Orders.Get(ctx, OrdersGetInput{
		TerritoryID: "territory",
		IDs: []string{
			"4516b2e1-43dc-49a8-8bfb-7190fa60df21",
			"0d56e7a3-c737-472e-bec9-e2f19d4865d3"},
	})

	c.Assert(err, qt.IsNil)
	c.Assert(len(o), qt.Equals, 2)
}

func TestOrdersAdd(t *testing.T) {
	setup()
	defer teardown()
	c := qt.New(t)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"requestId": "509900a5-392e-4d34-bcfe-90cc6bf3ad47"}`)
	})

	client, _ := New("api-key")
	client.baseURL, _ = url.Parse(server.URL)

	rID, err := client.Orders.Add(ctx, OrdersAddInput{
		TerritoryID:       "territory",
		Orders:            []Order{},
		Strict:            false,
		AcceptBadGeocodes: false,
	})

	c.Assert(err, qt.IsNil)
	c.Assert(rID, qt.Equals, "509900a5-392e-4d34-bcfe-90cc6bf3ad47")
}
