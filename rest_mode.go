// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to rest modes recorded by the Oura Ring
// Rest Mode API description: https://cloud.ouraring.com/v2/docs#tag/Rest-Mode-Period-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// RestModes stores a list of rest mode items along with a token which may be used to pull the next batch of RestMode items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_rest_mode_period_Documents_v2_usercollection_rest_mode_period_get
type RestModes struct {
	Items     []RestMode `json:"data"`
	NextToken string     `json:"next_token"`
}

// RestMode stores specifics for a recorded heart rate by the Oura Ring
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_rest_mode_period_Document_v2_usercollection_rest_mode_period__document_id__get
type RestMode struct {
	ID        string    `json:"id"`
	EndDay    Date      `json:"end_day"`
	EndTime   time.Time `json:"end_time"`
	Episodes  []Episode `json:"episodes"`
	StartDay  Date      `json:"start_day"`
	StartTime time.Time `json:"start_time"`
}

// Episode represents a single episode during a rest mode
type Episode struct {
	Tags      []string  `json:"tags"`
	Timestamp time.Time `json:"timestamp"`
}

type restModeBase RestMode
type restModesBase RestModes

// UnmarshalJSON is a helper function to convert a rest mode JSON from the API to the RestMode type.
func (rm *RestMode) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*rm), data); err != nil {
		return err
	}

	var restMode restModeBase
	err := json.Unmarshal(data, &restMode)
	if err != nil {
		return err
	}

	*rm = RestMode(restMode)
	return nil
}

// UnmarshalJSON is a helper function to convert a list of rest mode JSON from the API to the RestModes type.
func (rm *RestModes) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*rm), data); err != nil {
		return err
	}

	var restModes restModesBase
	err := json.Unmarshal(data, &restModes)
	if err != nil {
		return err
	}

	*rm = RestModes(restModes)
	return nil
}

// GetRestMode calls the Oura Ring API with a specific Rest Mode identifier and returns a RestMode object
func (c *Client) GetRestMode(restModeId string) (RestMode, error) {

	apiResponse, err := c.Getter(fmt.Sprintf(RestModeUrl+"/%s", restModeId), nil)

	if err != nil {
		return RestMode{},
			err
	}

	var restMode RestMode
	err = json.Unmarshal(*apiResponse, &restMode)
	if err != nil {
		return RestMode{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return restMode, nil
}

// GetRestModes accepts a start & end date and returns a RestModes object which will contain any RestMode
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// rest modes if the date range returns a large set.
func (c *Client) GetRestModes(startDate time.Time, endDate time.Time, nextToken *string) (RestModes, error) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, err := c.Getter(
		RestModeUrl,
		urlParameters,
	)

	if err != nil {
		return RestModes{},
			err
	}

	var restModes RestModes
	err = json.Unmarshal(*apiResponse, &restModes)
	if err != nil {
		return RestModes{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return restModes, nil
}
