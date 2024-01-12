package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type Stresses struct {
	Items     []Stress `json:"data"`
	NextToken string   `json:"next_token"`
}

type Stress struct {
	ID           string `json:"id"`
	Day          Date   `json:"day"`
	StressHigh   int64  `json:"stress_high"`
	RecoveryHigh int64  `json:"recovery_high"`
	DaySummary   string `json:"day_summary"`
}

type StressBase Stress
type StressesBase Stresses

func (sd *Stress) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*sd)
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

	var documentBase StressBase
	err = json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*sd = Stress(documentBase)
	return nil
}

func (sd *Stresses) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*sd)
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

	var documentBase StressesBase
	err = json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*sd = Stresses(documentBase)
	return nil
}

func (c *Client) GetStresses(startDate time.Time, endDate time.Time, nextToken *string) (Stresses, *OuraError) {
	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		StressUrl,
		urlParameters,
	)

	if ouraError != nil {
		return Stresses{},
			ouraError
	}

	var documents Stresses
	err := json.Unmarshal(*apiResponse, &documents)
	if err != nil {
		return Stresses{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return documents, nil
}

func (c *Client) GetStress(documentId string) (Stress, *OuraError) {
	apiResponse, ouraError := c.Getter(fmt.Sprintf(StressUrl+"/%s", documentId), nil)

	if ouraError != nil {
		return Stress{},
			ouraError
	}

	var stress Stress
	err := json.Unmarshal(*apiResponse, &stress)
	if err != nil {
		return Stress{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return stress, nil
}
