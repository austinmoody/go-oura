// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to heart rates recorded by the Oura Ring
// Heart Rate API description: https://cloud.ouraring.com/v2/docs#tag/Heart-Rate-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// HeartRates stores a list of heart rate items along with a token which may be used to pull the next batch of HeartRate items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_Heart_Rate_Documents_v2_usercollection_heartrate_get
type HeartRates struct {
	Items     []HeartRate `json:"data"`
	NextToken string      `json:"next_token"`
}

// HeartRate stores specifics for a recorded heart rate by the Oura Ring
type HeartRate struct {
	Bpm       int       `json:"bpm"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

type heartRatesBase HeartRates
type hearRateBase HeartRate

// UnmarshalJSON is a helper function to convert a heart rate JSON from the API to the HeartRate type.
func (hr *HeartRate) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*hr), data); err != nil {
		return err
	}

	var hrBase hearRateBase
	err := json.Unmarshal(data, &hrBase)
	if err != nil {
		return err
	}
	*hr = HeartRate(hrBase)

	return nil
}

// UnmarshalJSON is a helper function to convert a list of heart rate JSON from the API to the HeartRates type.
func (hr *HeartRates) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*hr), data); err != nil {
		return err
	}

	var hrBase heartRatesBase
	err := json.Unmarshal(data, &hrBase)
	if err != nil {
		return err
	}

	*hr = HeartRates(hrBase)

	return nil
}

// GetHeartRates accepts a start & end date and returns a HeartRates object which will contain any HeartRate
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// heart rates if the date range returns a large set.
func (c *Client) GetHeartRates(startDateTime time.Time, endDateTime time.Time, nextToken *string) (HeartRates, error) {

	urlParameters := url.Values{
		"start_datetime": []string{startDateTime.Format("2006-01-02T15:04:05-07:00")},
		"end_datetime":   []string{endDateTime.Format("2006-01-02T15:04:05-07:00")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, err := c.Getter(
		HeartRateUrl,
		urlParameters,
	)

	if err != nil {
		return HeartRates{}, err
	}

	var heartrates HeartRates
	err = json.Unmarshal(*apiResponse, &heartrates)
	if err != nil {
		return HeartRates{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return heartrates, nil
}
