package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

/*
https://cloud.ouraring.com/v2/docs#tag/Daily-Readiness-Routes
*/

type Readinesses struct {
	Items     []Readiness `json:"data"`
	NextToken *string     `json:"next_token"`
}

type Contributors struct {
	ActivityBalance     int `json:"activity_balance"`
	BodyTemperature     int `json:"body_temperature"`
	HrvBalance          int `json:"hrv_balance"`
	PreviousDayActivity int `json:"previous_day_activity"`
	PreviousNight       int `json:"previous_night"`
	RecoveryIndex       int `json:"recovery_index"`
	RestingHeartRate    int `json:"resting_heart_rate"`
	SleepBalance        int `json:"sleep_balance"`
}

type Readiness struct {
	Id                        string       `json:"id"`
	Contributors              Contributors `json:"contributors"`
	Day                       Date         `json:"day"`
	Score                     int          `json:"score"`
	TemperatureDeviation      float64      `json:"temperature_deviation"`
	TemperatureTrendDeviation float64      `json:"temperature_trend_deviation"`
	Timestamp                 time.Time    `json:"timestamp"`
}

type dailyReadinessDocumentBase Readiness
type dailyReadinessDocumentsBase Readinesses

// Custom UnmarshalJSON for Readiness and Readinesses
// Checks if all required fields are present in the JSON and returns
// an error if any are missing.
func (dr *Readinesses) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*dr)
	requiredFields := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		requiredFields = append(requiredFields, jsonTag)
	}

	for _, field := range requiredFields {
		if _, ok := rawMap[field]; !ok {
			return fmt.Errorf("required field %s not found", field)
		}
	}

	var documentBase dailyReadinessDocumentsBase
	err = json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*dr = Readinesses(documentBase)
	return nil
}

func (dr *Readiness) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*dr)
	requiredFields := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		requiredFields = append(requiredFields, jsonTag)
	}

	for _, field := range requiredFields {
		if _, ok := rawMap[field]; !ok {
			return fmt.Errorf("required field %s not found", field)
		}
	}

	var documentBase dailyReadinessDocumentBase
	err = json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*dr = Readiness(documentBase)
	return nil
}

func (c *Client) GetReadinesses(startDate time.Time, endDate time.Time, nextToken *string) (Readinesses, *OuraError) {

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
		return Readinesses{},
			ouraError
	}

	var readiness Readinesses
	err := json.Unmarshal(*apiResponse, &readiness)
	if err != nil {
		return Readinesses{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return readiness, nil
}

func (c *Client) GetReadiness(documentId string) (Readiness, *OuraError) {

	apiResponse, ouraError := c.Getter(fmt.Sprintf(ReadinessUrl+"/%s", documentId), nil)

	if ouraError != nil {
		return Readiness{},
			ouraError
	}

	var readiness Readiness
	err := json.Unmarshal(*apiResponse, &readiness)
	if err != nil {
		return Readiness{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return readiness, nil
}
