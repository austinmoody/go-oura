// Package go_oura provides a simple binding to the Oura Ring v2 API
package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// DailyReadinesses stores a list of daily DailyActivity items along with a token which may be used to pull the next batch of DailyActivity items from the API.
// https://cloud.ouraring.com/v2/docs#tag/Daily-Readiness-Routes
type DailyReadinesses struct {
	Items     []DailyReadiness `json:"data"`
	NextToken *string          `json:"next_token"`
}

// ReadinessContributors describes data points which contribute to the summary DailyReadiness score
type ReadinessContributors struct {
	ActivityBalance     int `json:"activity_balance"`
	BodyTemperature     int `json:"body_temperature"`
	HrvBalance          int `json:"hrv_balance"`
	PreviousDayActivity int `json:"previous_day_activity"`
	PreviousNight       int `json:"previous_night"`
	RecoveryIndex       int `json:"recovery_index"`
	RestingHeartRate    int `json:"resting_heart_rate"`
	SleepBalance        int `json:"sleep_balance"`
}

// DailyReadiness describes your daily readiness
type DailyReadiness struct {
	Id                        string                `json:"id"`
	Contributors              ReadinessContributors `json:"contributors"`
	Day                       Date                  `json:"day"`
	Score                     int                   `json:"score"`
	TemperatureDeviation      float64               `json:"temperature_deviation"`
	TemperatureTrendDeviation float64               `json:"temperature_trend_deviation"`
	Timestamp                 time.Time             `json:"timestamp"`
}

type dailyReadinessDocumentBase DailyReadiness
type dailyReadinessDocumentsBase DailyReadinesses

// UnmarshalJSON is a helper function to convert daily readinesses JSON from the API to the DailyReadinesses type.
func (dr *DailyReadinesses) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*dr), data); err != nil {
		return err
	}

	var documentBase dailyReadinessDocumentsBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*dr = DailyReadinesses(documentBase)
	return nil
}

// UnmarshalJSON is a helper function to convert a daily readiness JSON from the API to the DailyReadiness type.
func (dr *DailyReadiness) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*dr), data); err != nil {
		return err
	}

	var documentBase dailyReadinessDocumentBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*dr = DailyReadiness(documentBase)
	return nil
}

// GetReadinesses accepts a start & end date and returns a DailyReadinesses object which will contain any DailyReadiness
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// activities if the date range returns a large set.
func (c *Client) GetReadinesses(startDate time.Time, endDate time.Time, nextToken *string) (DailyReadinesses, *OuraError) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		ReadinessUrl,
		urlParameters,
	)

	if ouraError != nil {
		return DailyReadinesses{},
			ouraError
	}

	var readiness DailyReadinesses
	err := json.Unmarshal(*apiResponse, &readiness)
	if err != nil {
		return DailyReadinesses{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return readiness, nil
}

// GetReadiness accepts a single Daily Readiness ID and returns a DailyReadiness object.
func (c *Client) GetReadiness(dailyReadinessId string) (DailyReadiness, *OuraError) {

	apiResponse, ouraError := c.Getter(fmt.Sprintf(ReadinessUrl+"/%s", dailyReadinessId), nil)

	if ouraError != nil {
		return DailyReadiness{},
			ouraError
	}

	var readiness DailyReadiness
	err := json.Unmarshal(*apiResponse, &readiness)
	if err != nil {
		return DailyReadiness{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return readiness, nil
}
