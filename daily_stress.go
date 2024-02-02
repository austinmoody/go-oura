// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to Daily Stress recorded by the Oura Ring
// Daily Sleep API description: https://cloud.ouraring.com/v2/docs#tag/Daily-Stress-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// DailyStresses stores a list of daily stress items along with a token which may be used to pull the next batch of DailyStress items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_daily_stress_Documents_v2_usercollection_daily_stress_get
type DailyStresses struct {
	Items     []DailyStress `json:"data"`
	NextToken string        `json:"next_token"`
}

// DailyStress describes daily stress summary values
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_daily_stress_Document_v2_usercollection_daily_stress__document_id__get
type DailyStress struct {
	ID           string `json:"id"`
	Day          Date   `json:"day"`
	StressHigh   int64  `json:"stress_high"`
	RecoveryHigh int64  `json:"recovery_high"`
	DaySummary   string `json:"day_summary"`
}

type StressBase DailyStress
type StressesBase DailyStresses

// UnmarshalJSON is a helper function to convert a daily stress JSON from the API to the DailyStress type.
func (sd *DailyStress) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*sd), data); err != nil {
		return err
	}

	var documentBase StressBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*sd = DailyStress(documentBase)
	return nil
}

// UnmarshalJSON is a helper function to convert daily stresses JSON from the API to the DailyStresses type.
func (sd *DailyStresses) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*sd), data); err != nil {
		return err
	}

	var documentBase StressesBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*sd = DailyStresses(documentBase)
	return nil
}

// GetStresses accepts a start & end date and returns a DailyStresses object which will contain any DailyStress
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// stresses if the date range returns a large set.
func (c *Client) GetStresses(startDate time.Time, endDate time.Time, nextToken *string) (DailyStresses, error) {
	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, err := c.Getter(
		StressUrl,
		urlParameters,
	)

	if err != nil {
		return DailyStresses{},
			err
	}

	var documents DailyStresses
	err = json.Unmarshal(*apiResponse, &documents)
	if err != nil {
		return DailyStresses{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return documents, nil
}

// GetStress accepts a single Daily Stress ID and returns a DailyStress object.
func (c *Client) GetStress(stressId string) (DailyStress, error) {
	apiResponse, err := c.Getter(fmt.Sprintf(StressUrl+"/%s", stressId), nil)

	if err != nil {
		return DailyStress{},
			err
	}

	var stress DailyStress
	err = json.Unmarshal(*apiResponse, &stress)
	if err != nil {
		return DailyStress{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return stress, nil
}
