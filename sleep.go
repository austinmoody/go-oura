// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to Sleep recorded by the Oura Ring
// https://cloud.ouraring.com/v2/docs#tag/Sleep-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// Sleeps stores a list of sleep items along with a token which may be used to pull the next batch of Sleep items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_sleep_Documents_v2_usercollection_sleep_get
type Sleeps struct {
	Items     []Sleep `json:"data"`
	NextToken string  `json:"next_token"`
}

// Sleep describes a single sleep session
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_sleep_Document_v2_usercollection_sleep__document_id__get
type Sleep struct {
	ID                    string         `json:"id"`
	AverageBreath         float64        `json:"average_breath"`
	AverageHeartRate      float64        `json:"average_heart_rate"`
	AverageHrv            int            `json:"average_hrv"`
	AwakeTime             int            `json:"awake_time"`
	BedtimeEnd            time.Time      `json:"bedtime_end"`
	BedtimeStart          time.Time      `json:"bedtime_start"`
	Day                   Date           `json:"day"`
	DeepSleepDuration     int            `json:"deep_sleep_duration"`
	Efficiency            int            `json:"efficiency"`
	HeartRate             IntervalItems  `json:"heart_rate"`
	Hrv                   IntervalItems  `json:"hrv"`
	Latency               int            `json:"latency"`
	LightSleepDuration    int            `json:"light_sleep_duration"`
	LowBatteryAlert       bool           `json:"low_battery_alert"`
	LowestHeartRate       int            `json:"lowest_heart_rate"`
	Movement30Sec         string         `json:"movement_30_sec"`
	Period                int            `json:"period"`
	Readiness             SleepReadiness `json:"readiness"`
	ReadinessScoreDelta   int            `json:"readiness_score_delta"`
	RemSleepDuration      int            `json:"rem_sleep_duration"`
	RestlessPeriods       int            `json:"restless_periods"`
	SleepPhase5Min        string         `json:"sleep_phase_5_min"`
	SleepScoreDelta       int            `json:"sleep_score_delta"`
	SleepAlgorithmVersion string         `json:"sleep_algorithm_version"`
	TimeInBed             int            `json:"time_in_bed"`
	TotalSleepDuration    int            `json:"total_sleep_duration"`
	Type                  string         `json:"type"`
}

type SleepReadiness struct {
	Contributors              Contributors `json:"contributors"`
	Score                     int          `json:"score"`
	TemperatureDeviation      float64      `json:"temperature_deviation"`
	TemperatureTrendDeviation float64      `json:"temperature_trend_deviation"`
}

type sleepDocumentBase Sleep
type sleepDocumentsBase Sleeps

// UnmarshalJSON is a helper function to convert a sleep JSON from the API to the Sleep type.
func (s *Sleep) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*s), data); err != nil {
		return err
	}

	var documentBase sleepDocumentBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*s = Sleep(documentBase)
	return nil
}

// UnmarshalJSON is a helper function to convert multiple sleep JSON from the API to the Sleeps type.
func (s *Sleeps) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*s), data); err != nil {
		return err
	}

	var documentBase sleepDocumentsBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*s = Sleeps(documentBase)
	return nil
}

// GetSleeps accepts a start & end date and returns a Sleeps object which will contain any Sleep
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// sleep if the date range returns a large set.
func (c *Client) GetSleeps(startDate time.Time, endDate time.Time, nextToken *string) (Sleeps, *OuraError) {

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
		return Sleeps{},
			ouraError
	}

	var documents Sleeps
	err := json.Unmarshal(*apiResponse, &documents)
	if err != nil {
		return Sleeps{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return documents, nil
}

// GetSleep accepts a single sleep ID and returns a Sleep object.
func (c *Client) GetSleep(sleepId string) (Sleep, *OuraError) {

	apiResponse, ouraError := c.Getter(
		fmt.Sprintf(SleepUrl+"/%s", sleepId),
		nil,
	)

	if ouraError != nil {
		return Sleep{},
			ouraError
	}

	var sleep Sleep
	err := json.Unmarshal(*apiResponse, &sleep)
	if err != nil {
		return Sleep{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return sleep, nil
}
