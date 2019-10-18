package workwave

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	callbackPath = "/api/v1/callback"
)

// CallbackService is an interface to the callback configuration for the WorkWave API.
type CallbackService interface {
	Get(context.Context, Callback) (Callback, error)
	Set(context.Context, Callback) (Callback, error)
	Delete(context.Context, Callback) (Callback, error)
}

type callbackService struct {
	client *Client
}

// Callback is a multi-usage type that contains both fields to create and receive
// replies back from the WorkWave callback API.
type Callback struct {
	URL         string `json:"url,omitempty"`
	PreviousURL string `json:"previousUrl,omitempty"`
	// SignaturePassword is optional and is only used for Set
	// For usage, see reference: https://wwrm.workwave.com/api/#set-callback-url
	SignaturePassword string `json:"signaturePassword,omitempty"`
	// Test is optional. If set to true, WorkWave will test the given callback URL
	// synchronously at the time the set callback call is made. ErrorCode and
	// ErrorMessage will be populated in the response if the test has been requested
	// and fails.
	Test bool `json:"test,omitempty"`
	// Headers are optional. If provided, the given headers will be added to each
	// callback POST from WorkWave. The concatenated string of all keys and values in
	// headers must not exceed 256k characters.
	Headers map[string]string `json:"headers,omitempty"`
	// ErrorCode and ErrorMessage are only present if a test is requested
	// and the check from WorkWave fails.
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func (svc *callbackService) Get(ctx context.Context, c Callback) (Callback, error) {
	callback := Callback{}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, callbackPath, c)
	if err != nil {
		return callback, errors.Wrap(err, "failed to create callback get request")
	}

	if _, err := svc.client.Do(ctx, req, &callback); err != nil {
		return callback, nil
	}
	return callback, nil
}

func (svc *callbackService) Set(ctx context.Context, c Callback) (Callback, error) {
	callback := Callback{}
	req, err := svc.client.NewRequest(ctx, http.MethodPost, callbackPath, c)
	if err != nil {
		return callback, errors.Wrap(err, "failed to create callback set request")
	}

	if _, err := svc.client.Do(ctx, req, &callback); err != nil {
		return callback, nil
	}
	if callback.ErrorCode != 0 {
		return callback, fmt.Errorf("failed to set callback: %s", callback.ErrorMessage)
	}
	return callback, nil
}

func (svc *callbackService) Delete(ctx context.Context, c Callback) (Callback, error) {
	callback := Callback{}
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, callbackPath, c)
	if err != nil {
		return callback, errors.Wrap(err, "failed to create callback delete request")
	}

	if _, err := svc.client.Do(ctx, req, &callback); err != nil {
		return callback, nil
	}
	return callback, nil
}
