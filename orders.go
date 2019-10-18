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
	List(context.Context, OrderListInput) ([]Order, error)
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

// OrderListInput is used to populate a call List Orders on the WorkWave API.
type OrderListInput struct {
	TerritoryID string
	Include     string
	EligibleOn  string
	AssignedOn  string
}

type orderListResponse struct {
	Orders map[string]interface{}
}

// List retrieves the details for a Story.
func (svc *ordersService) List(ctx context.Context, i OrderListInput) ([]Order, error) {
	u := fmt.Sprintf(ordersBasePath, i.TerritoryID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create order list request")
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

	olr := &orderListResponse{}
	svc.client.Do(ctx, req, olr)

	// To be continued …
	return nil, nil
}
