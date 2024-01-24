// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to sleep time, which according to the Oura Ring documentation are:
// "Recommendations for the optimal bedtime window that is calculated based on sleep data."
// https://cloud.ouraring.com/v2/docs#tag/Sleep-Time-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// SleepTimes stores a list of sleep time items along with a token which may be used to pull the next batch of SleepTime items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_sleep_time_Documents_v2_usercollection_sleep_time_get
type SleepTimes struct {
	Items     []SleepTime `json:"data"`
	NextToken string      `json:"next_token"`
}

// SleepTime stores specifics for a single day optimal bedtime window
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_sleep_time_Document_v2_usercollection_sleep_time__document_id__get
type SleepTime struct {
	ID             string          `json:"id"`
	Day            Date            `json:"day"`
	OptimalBedtime *OptimalBedtime `json:"optimal_bedtime"`
	Recommendation string          `json:"recommendation"`
	Status         string          `json:"status"`
}

type OptimalBedtime struct {
	DayTz       int `json:"day_tz"`
	EndOffset   int `json:"end_offset"`
	StartOffset int `json:"start_offset"`
}

type sleepTimeBase SleepTime
type sleepTimesBase SleepTimes

// UnmarshalJSON is a helper function to convert a sleep time JSON from the API to the SleepTime type.
func (st *SleepTime) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*st), data); err != nil {
		return err
	}

	var sleepTime sleepTimeBase
	err := json.Unmarshal(data, &sleepTime)
	if err != nil {
		return err
	}

	*st = SleepTime(sleepTime)
	return nil
}

// UnmarshalJSON is a helper function to convert multiple sleep times JSON from the API to the SleepTimes type.
func (st *SleepTimes) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*st), data); err != nil {
		return err
	}

	var sleepTimes sleepTimesBase
	err := json.Unmarshal(data, &sleepTimes)
	if err != nil {
		return err
	}

	*st = SleepTimes(sleepTimes)
	return nil
}

// GetSleepTimes accepts a start & end date and returns a SleepTimes object which will contain any SleepTime
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// sleep times if the date range returns a large set.
func (c *Client) GetSleepTimes(startDate time.Time, endDate time.Time, nextToken *string) (SleepTimes, *OuraError) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		SleepTimeUrl,
		urlParameters,
	)

	if ouraError != nil {
		return SleepTimes{},
			ouraError
	}

	var sleepTimes SleepTimes
	err := json.Unmarshal(*apiResponse, &sleepTimes)
	if err != nil {
		return SleepTimes{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return sleepTimes, nil
}

// GetSleepTime calls the Oura Ring API with a specific sleep identifer and returns a SleepTime object
func (c *Client) GetSleepTime(sleepTimeId string) (SleepTime, *OuraError) {

	apiResponse, ouraError := c.Getter(
		fmt.Sprintf(SleepTimeUrl+"/%s", sleepTimeId),
		nil,
	)

	if ouraError != nil {
		return SleepTime{},
			ouraError
	}

	var sleepTime SleepTime
	err := json.Unmarshal(*apiResponse, &sleepTime)
	if err != nil {
		return SleepTime{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return sleepTime, nil
}
