package workwave

// Route represents a route in WorkWave which is associated with a date,
// vehicle, driver and steps to complete order deliveries.
type Route struct {
	ID        string      `json:"id,omitempty"`
	Revision  int         `json:"revision,omitempty"`
	Date      string      `json:"date,omitempty"` // in the format yyyyMMdd
	Steps     []RouteStep `json:"steps,omitempty"`
	DriverID  string      `json:"driverId,omitempty"`
	VehicleID string      `json:"vehicleId,omitempty"`
}

// RouteStep is one step along a delivery route and include departure,
// a number of deliveries, and arrival.
type RouteStep struct {
	Type         string        `json:"type,omitempty"` // One of: departure, arrival, pickup. delivery, brk
	OrderID      string        `json:"orderId,omitempty"`
	ArrivalSec   int           `json:"arrivalSec,omitempty"`
	StartSec     int           `json:"startSec,omitempty"`
	EndSec       int           `json:"endSec,omitempty"`
	DisplayLabel string        `json:"displayLabel,omitempty"`
	TrackingData *TrackingData `json:"trackingData,omitempty"`
}

// TrackingData provides location, timing and status for a route step.
type TrackingData struct {
	Status    string `json:"status,omitempty"`
	StatusSec int    `json:"statusSec,omitempty"`
}

// RouteList maps the response from calling the "List Approved Routes"
// endpoint in WorkWave.
type RouteList struct {
	Routes   map[string]*Route   `json:"routes,omitempty"`
	Vehicles map[string]*Vehicle `json:"vehicles,omitempty"`
	Drivers  map[string]*Driver  `json:"drivers,omitempty"`
}
