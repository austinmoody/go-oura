// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to Daily Readiness recorded by the Oura Ring
// Daily Readiness API description: https://cloud.ouraring.com/v2/docs#tag/Daily-Readiness-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// DailyReadinesses stores a list of daily DailyActivity items along with a token which may be used to pull the next batch of DailyActivity items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_daily_readiness_Documents_v2_usercollection_daily_readiness_get
type DailyReadinesses struct {
	Items     []DailyReadiness `json:"data"`
	NextToken *string          `json:"next_token"`
}

// ReadinessContributors describes data points which contribute to the summary DailyReadiness score
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_daily_readiness_Document_v2_usercollection_daily_readiness__document_id__get
type ReadinessContributors Contributors

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
func (c *Client) GetReadinesses(startDate time.Time, endDate time.Time, nextToken *string) (DailyReadinesses, error) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, err := c.Getter(
		ReadinessUrl,
		urlParameters,
	)

	if err != nil {
		return DailyReadinesses{},
			err
	}

	var readiness DailyReadinesses
	err = json.Unmarshal(*apiResponse, &readiness)
	if err != nil {
		return DailyReadinesses{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return readiness, nil
}

// GetReadiness accepts a single Daily Readiness ID and returns a DailyReadiness object.
func (c *Client) GetReadiness(dailyReadinessId string) (DailyReadiness, error) {

	apiResponse, err := c.Getter(fmt.Sprintf(ReadinessUrl+"/%s", dailyReadinessId), nil)

	if err != nil {
		return DailyReadiness{},
			err
	}

	var readiness DailyReadiness
	err = json.Unmarshal(*apiResponse, &readiness)
	if err != nil {
		return DailyReadiness{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return readiness, nil
}
