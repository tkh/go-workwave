package workwave

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	routesBasePath = "/api/v1/territories/%s"
	toaRoutesPath  = routesBasePath + "/toa/routes"
	// approvedRoutesPath is for the v1 API
	// v2 also exists: https://wwrm.workwave.com/api/#approved-plans-api-v2
	approvedRoutesPath = routesBasePath + "/approved/routes"
)

// RoutesService is an interface to routes in the WorkWave API.
// Routes are available both through the Time of Arrival and Approved Plans APIs.
type RoutesService interface {
	ListCurrent(context.Context, RoutesListCurrentInput) ([]Route, error)
	ListApproved(context.Context, RoutesListApprovedInput) ([]Route, error)
}

type routesService struct {
	client *Client
}

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

// RoutesListCurrentInput is used to populate a call to List Current Routes on the
// WorkWave API.
type RoutesListCurrentInput struct {
	TerritoryID string `json:"-"`
	Date        string `json:"date"`
	Vehicle     string `json:"vehicle"`
}

type routesListResponse struct {
	Routes map[string]Route `json:"routes"`
}

// ListCurrent lists current, live Routes, optionally filtering by date
// and/or vehicleId.
func (svc *routesService) ListCurrent(ctx context.Context, i RoutesListCurrentInput) ([]Route, error) {
	u := fmt.Sprintf(toaRoutesPath, i.TerritoryID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create current route list request")
	}

	// Query params
	q := req.URL.Query()
	if i.Date != "" {
		q.Add("date", i.Date)
	}
	if i.Vehicle != "" {
		q.Add("vehicle", i.Vehicle)
	}

	rlr := &routesListResponse{}
	if _, err := svc.client.Do(ctx, req, rlr); err != nil {
		return nil, err
	}

	var routes []Route
	for _, route := range rlr.Routes {
		routes = append(routes, route)
	}

	return routes, nil
}

// RoutesListApprovedInput
type RoutesListApprovedInput struct {
	TerritoryID string `json:"-"`
	Date        string `json:"date"`
}

// ListApproved lists approved planned routes.
func (svc *routesService) ListApproved(ctx context.Context, i RoutesListApprovedInput) ([]Route, error) {
	u := fmt.Sprintf(approvedRoutesPath, i.TerritoryID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create approved route list request")
	}

	// Query params
	q := req.URL.Query()
	if i.Date != "" {
		q.Add("date", i.Date)
	}

	rlr := &routesListResponse{}
	if _, err := svc.client.Do(ctx, req, rlr); err != nil {
		return nil, err
	}

	var routes []Route
	for _, route := range rlr.Routes {
		routes = append(routes, route)
	}

	return routes, nil
}
