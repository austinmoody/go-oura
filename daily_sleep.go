// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to Daily Sleep recorded by the Oura Ring
// Daily Sleep API description: https://cloud.ouraring.com/v2/docs#tag/Daily-Sleep-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// DailySleeps stores a list of daily sleep items along with a token which may be used to pull the next batch of DailySleep items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_daily_sleep_Documents_v2_usercollection_daily_sleep_get
type DailySleeps struct {
	Items     []DailySleep `json:"data"`
	NextToken string       `json:"next_token"`
}

// DailySleep describes a single sleep session
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_daily_sleep_Document_v2_usercollection_daily_sleep__document_id__get
type DailySleep struct {
	ID           string            `json:"id"`
	Contributors SleepContributors `json:"contributors"`
	Day          Date              `json:"day"`
	Score        int64             `json:"score"`
	Timestamp    time.Time         `json:"timestamp"`
}

// SleepContributors describes data points which contribute to the DailySleep score
type SleepContributors struct {
	DeepSleep   int64 `json:"deep_sleep"`
	Efficiency  int64 `json:"efficiency"`
	Latency     int64 `json:"latency"`
	RemSleep    int64 `json:"rem_sleep"`
	Restfulness int64 `json:"restfulness"`
	Timing      int64 `json:"timing"`
	TotalSleep  int64 `json:"total_sleep"`
}

type dailySleepDocumentBase DailySleep
type dailySleepDocumentsBase DailySleeps

// UnmarshalJSON is a helper function to convert a daily sleep JSON from the API to the DailySleep type.
func (sd *DailySleep) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*sd), data); err != nil {
		return err
	}

	var documentBase dailySleepDocumentBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*sd = DailySleep(documentBase)
	return nil
}

// UnmarshalJSON is a helper function to convert multiple daily sleep JSON from the API to the DailySleeps type.
func (sd *DailySleeps) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*sd), data); err != nil {
		return err
	}

	var documentBase dailySleepDocumentsBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*sd = DailySleeps(documentBase)
	return nil
}

// GetSleeps accepts a start & end date and returns a DailySleeps object which will contain any DailySleep
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// activities if the date range returns a large set.
func (c *Client) GetSleeps(startDate time.Time, endDate time.Time, nextToken *string) (DailySleeps, *OuraError) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		SleepUrl,
		urlParameters,
	)

	if ouraError != nil {
		return DailySleeps{},
			ouraError
	}

	var documents DailySleeps
	err := json.Unmarshal(*apiResponse, &documents)
	if err != nil {
		return DailySleeps{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return documents, nil
}

// GetSleep accepts a single daily sleep ID and returns a DailySleep object.
func (c *Client) GetSleep(dailySleepId string) (DailySleep, *OuraError) {

	apiResponse, ouraError := c.Getter(fmt.Sprintf(SleepUrl+"/%s", dailySleepId), nil)

	if ouraError != nil {
		return DailySleep{},
			ouraError
	}

	var sleep DailySleep
	err := json.Unmarshal(*apiResponse, &sleep)
	if err != nil {
		return DailySleep{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return sleep, nil
}
