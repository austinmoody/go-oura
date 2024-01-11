package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type HeartRates struct {
	Items     []HeartRate `json:"data"`
	NextToken string      `json:"next_token"`
}

type HeartRate struct {
	Bpm       int       `json:"bpm"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

type heartRatesBase HeartRates
type hearRateBase HeartRate

func (hr *HeartRate) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*hr)
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

	var hrBase hearRateBase
	err = json.Unmarshal(data, &hrBase)
	if err != nil {
		return err
	}
	*hr = HeartRate(hrBase)

	return nil
}

func (hr *HeartRates) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*hr)
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

	var hrBase heartRatesBase
	err = json.Unmarshal(data, &hrBase)
	if err != nil {
		return err
	}

	*hr = HeartRates(hrBase)

	return nil
}

func (c *Client) GetHeartRates(startDateTime time.Time, endDateTime time.Time) (HeartRates, *OuraError) {
	apiResponse, ouraError := c.Getter(
		HeartRateUrl,
		url.Values{
			"start_datetime": []string{startDateTime.Format("2006-01-02T15:04:05-07:00")},
			"end_datetime":   []string{endDateTime.Format("2006-01-02T15:04:05-07:00")},
		},
	)

	if ouraError != nil {
		return HeartRates{}, ouraError
	}

	var heartrates HeartRates
	err := json.Unmarshal(*apiResponse, &heartrates)
	if err != nil {
		return HeartRates{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return heartrates, nil
}
