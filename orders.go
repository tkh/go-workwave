package workwave

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
