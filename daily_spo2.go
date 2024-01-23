// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to Daily SpO2 recorded by the Oura Ring
// SpO2 = Blood Oxygenation
// Only available starting with Gen 3 Oura Rings
// Daily SpO2 API description: https://cloud.ouraring.com/v2/docs#tag/Daily-Spo2-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// DailySpo2Readings stores a list of daily spo2 readings along with a token which may be used to pull the next batch of DailySpo2Reading items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_daily_spo2_Documents_v2_usercollection_daily_spo2_get
type DailySpo2Readings struct {
	Items     []DailySpo2Reading `json:"data"`
	NextToken *string            `json:"next_token"`
}

// DailySpo2Reading include daily SpO2 average.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_daily_spo2_Document_v2_usercollection_daily_spo2__document_id__get
type DailySpo2Reading struct {
	ID         string         `json:"id"`
	Day        Date           `json:"day"`
	Percentage Spo2Percentage `json:"spo2_percentage"`
}

// Spo2Percentage is a simple type that currently only holds the Average SpO2 for a reading
type Spo2Percentage struct {
	Average float64 `json:"average"`
}

type spo2Base DailySpo2Reading
type spo2sBase DailySpo2Readings

// UnmarshalJSON is a helper function to convert a single SpO2 reading JSON from the API to the DailySpo2Reading type.
func (s *DailySpo2Reading) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*s), data); err != nil {
		return err
	}

	var spo2 spo2Base
	err := json.Unmarshal(data, &spo2)
	if err != nil {
		return err
	}

	*s = DailySpo2Reading(spo2)
	return nil
}

// UnmarshalJSON is a helper function to convert daily SpO2 readings JSON from the API to the DailySpo2Readings type.
func (s *DailySpo2Readings) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*s), data); err != nil {
		return err
	}

	var spo2 spo2sBase
	err := json.Unmarshal(data, &spo2)
	if err != nil {
		return err
	}

	*s = DailySpo2Readings(spo2)

	return nil
}

// GetSpo2Readings accepts a start & end date and returns a DailySpo2Readings object which will contain any DailySpo2Reading
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// SpO2 readings if the date range returns a large set.
func (c *Client) GetSpo2Readings(startDate time.Time, endDate time.Time, nextToken *string) (DailySpo2Readings, *OuraError) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		Spo2Url,
		urlParameters,
	)

	if ouraError != nil {
		return DailySpo2Readings{}, ouraError
	}

	var readings DailySpo2Readings
	err := json.Unmarshal(*apiResponse, &readings)
	if err != nil {
		return DailySpo2Readings{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return readings, nil
}

// GetSpo2Reading accepts a single SpO2 Reading ID and returns a DailySpo2Reading object.
func (c *Client) GetSpo2Reading(spo2ReadingId string) (DailySpo2Reading, *OuraError) {
	apiResponse, ouraError := c.Getter(
		fmt.Sprintf("%s/%s", Spo2Url, spo2ReadingId),
		nil,
	)

	if ouraError != nil {
		return DailySpo2Reading{}, ouraError
	}

	var reading DailySpo2Reading
	err := json.Unmarshal(*apiResponse, &reading)
	if err != nil {
		return DailySpo2Reading{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return reading, nil
}
