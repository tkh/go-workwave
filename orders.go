package workwave

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	ordersBasePath = "/api/v1/territories/%s/orders"
)

// OrdersService is an interface to orders in the WorkWave API.
type OrdersService interface {
	List(context.Context, OrdersListInput) ([]Order, error)
	Get(context.Context, OrdersGetInput) ([]Order, error)
}

type ordersService struct {
	client *Client
}

// Order represents an Order in the WorkWave API
// This structure can be used as input for order calls by omitting ID.
type Order struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	Eligibility    Eligibility    `json:"eligibility,omitempty"`
	ForceVehicleID interface{}    `json:"forceVehicleId,omitempty"`
	Priority       int            `json:"priority,omitempty"`
	Loads          map[string]int `json:"loads,omitempty"`
	Pickup         *OrderStep     `json:"pickup,omitempty"`
	Delivery       *OrderStep     `json:"delivery,omitempty"`
	IsService      bool           `json:"isService,omitempty"`
}

// Eligibility represents Eligibility for an Order in the WorkWave API.
type Eligibility struct {
	Type    string   `json:"type,omitempty"`    // One of: on, by, any
	ByDate  string   `json:"byDate,omitempty"`  // Used when type = by, format: yyyyMMdd
	OnDates []string `json:"onDates,omitempty"` // Used when type = on, format: yyyyMMdd
}

// Location represents a Location in the WorkWave API.
type Location struct {
	Address string  `json:"address,omitempty"`
	LatLng  *[2]int `json:"latLng,omitempty"` // ie, {33817872, -87266893}
	Status  string  `json:"status,omitempty"`
}

// TimeWindow represents a time window in the WorkWave API.
type TimeWindow struct {
	StartSec int `json:"startSec,omitempty"`
	EndSec   int `json:"endSec,omitempty"`
}

// OrderStep represents an OrderStep within an Order in the WorkWave API.
// An OrderStep can be `pickup` or `delivery`.
type OrderStep struct {
	DepotID              string                `json:"depotId,omitempty"`
	Location             Location              `json:"location,omitempty"`
	TimeWindows          []TimeWindow          `json:"timeWindows,omitempty"`
	TimeWindowExceptions map[string]TimeWindow `json:"timeWindowExceptions,omitempty"`
	Notes                string                `json:"notes,omitempty"`
	ServiceTimeSec       int                   `json:"serviceTimeSec,omitempty"`
	TagsIn               []string              `json:"tagsIn,omitempty"`
	TagsOut              []string              `json:"tagsOut,omitempty"`
	CustomFields         map[string]string     `json:"customFields,omitempty"`
}

type ordersResponse struct {
	Orders map[string]Order `json:"orders"`
}

// OrdersListInput is used to populate a call to List Orders on the WorkWave API.
type OrdersListInput struct {
	TerritoryID string
	Include     string
	EligibleOn  string
	AssignedOn  string
}

// List retrieves the orders matching the filters provided in the given OrderListInput.
func (svc *ordersService) List(ctx context.Context, i OrdersListInput) ([]Order, error) {
	u := fmt.Sprintf(ordersBasePath, i.TerritoryID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create orders list request")
	}

	// Query params
	q := req.URL.Query()
	if i.Include != "" {
		q.Add("include", i.Include)
	}
	if i.EligibleOn != "" {
		q.Add("eligibleOn", i.EligibleOn)
	}
	if i.AssignedOn != "" {
		q.Add("assignedOn", i.AssignedOn)
	}
	req.URL.RawQuery = q.Encode()

	olr := &ordersResponse{}
	if _, err := svc.client.Do(ctx, req, olr); err != nil {
		return nil, err
	}

	var orders []Order
	for _, order := range olr.Orders {
		orders = append(orders, order)
	}

	return orders, nil
}

// OrdersGetInput is used to populate a call Get Orders on the WorkWave API.
type OrdersGetInput struct {
	TerritoryID string
	IDs         []string
}

// Get orders for the given IDs.
func (svc *ordersService) Get(ctx context.Context, i OrdersGetInput) ([]Order, error) {
	u := fmt.Sprintf(ordersBasePath, i.TerritoryID)
	b := struct {
		IDs []string `json:"ids"`
	}{
		IDs: i.IDs,
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, u, b)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create orders get request")
	}

	olr := &ordersResponse{}
	if _, err := svc.client.Do(ctx, req, olr); err != nil {
		return nil, err
	}

	var orders []Order
	for _, order := range olr.Orders {
		orders = append(orders, order)
	}

	return orders, nil
}
