package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type Spo2Readings struct {
	Items     []Spo2Reading `json:"data"`
	NextToken *string       `json:"next_token"`
}

type Spo2Reading struct {
	ID         string         `json:"id"`
	Day        Date           `json:"day"`
	Percentage Spo2Percentage `json:"spo2_percentage"`
}

type Spo2Percentage struct {
	Average float64 `json:"average"`
}

type spo2Base Spo2Reading
type spo2sBase Spo2Readings

func (s *Spo2Reading) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*s)
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

	var spo2 spo2Base
	err = json.Unmarshal(data, &spo2)
	if err != nil {
		return err
	}

	*s = Spo2Reading(spo2)
	return nil
}

func (s *Spo2Readings) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*s)
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

	var spo2 spo2sBase
	err = json.Unmarshal(data, &spo2)
	if err != nil {
		return err
	}

	*s = Spo2Readings(spo2)

	return nil
}

func (c *Client) GetSpo2Readings(startDate time.Time, endDate time.Time) (Spo2Readings, *OuraError) {

	apiResponse, ouraError := c.Getter(
		Spo2Url,
		url.Values{
			"start_date": []string{startDate.Format("2006-01-02")},
			"end_date":   []string{endDate.Format("2006-01-02")},
		},
	)

	if ouraError != nil {
		return Spo2Readings{}, ouraError
	}

	var readings Spo2Readings
	err := json.Unmarshal(*apiResponse, &readings)
	if err != nil {
		return Spo2Readings{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return readings, nil
}

func (c *Client) GetSpo2Reading(id string) (Spo2Reading, *OuraError) {
	apiResponse, ouraError := c.Getter(
		fmt.Sprintf("%s/%s", Spo2Url, id),
		nil,
	)

	if ouraError != nil {
		return Spo2Reading{}, ouraError
	}

	var reading Spo2Reading
	err := json.Unmarshal(*apiResponse, &reading)
	if err != nil {
		return Spo2Reading{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return reading, nil
}
