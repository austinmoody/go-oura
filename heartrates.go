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
	if err := checkJSONFields(reflect.TypeOf(*hr), data); err != nil {
		return err
	}

	var hrBase hearRateBase
	err := json.Unmarshal(data, &hrBase)
	if err != nil {
		return err
	}
	*hr = HeartRate(hrBase)

	return nil
}

func (hr *HeartRates) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*hr), data); err != nil {
		return err
	}

	var hrBase heartRatesBase
	err := json.Unmarshal(data, &hrBase)
	if err != nil {
		return err
	}

	*hr = HeartRates(hrBase)

	return nil
}

func (c *Client) GetHeartRates(startDateTime time.Time, endDateTime time.Time, nextToken *string) (HeartRates, *OuraError) {

	urlParameters := url.Values{
		"start_datetime": []string{startDateTime.Format("2006-01-02T15:04:05-07:00")},
		"end_datetime":   []string{endDateTime.Format("2006-01-02T15:04:05-07:00")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		HeartRateUrl,
		urlParameters,
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
