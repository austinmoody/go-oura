package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type SleepTimes struct {
	Items     []SleepTime `json:"data"`
	NextToken string      `json:"next_token"`
}

type SleepTime struct {
	ID             string          `json:"id"`
	Day            Date            `json:"day"`
	OptimalBedtime *OptimalBedtime `json:"optimal_bedtime"`
	Recommendation string          `json:"recommendation"`
	Status         string          `json:"status"`
}

type OptimalBedtime struct {
	DayTz       int `json:"day_tz"`
	EndOffset   int `json:"end_offset"`
	StartOffset int `json:"start_offset"`
}

type sleepTimeBase SleepTime
type sleepTimesBase SleepTimes

func (st *SleepTime) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*st), data); err != nil {
		return err
	}

	var sleepTime sleepTimeBase
	err := json.Unmarshal(data, &sleepTime)
	if err != nil {
		return err
	}

	*st = SleepTime(sleepTime)
	return nil
}

func (st *SleepTimes) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*st), data); err != nil {
		return err
	}

	var sleepTimes sleepTimesBase
	err := json.Unmarshal(data, &sleepTimes)
	if err != nil {
		return err
	}

	*st = SleepTimes(sleepTimes)
	return nil
}

func (c *Client) GetSleepTimes(startDate time.Time, endDate time.Time, nextToken *string) (SleepTimes, *OuraError) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		SleepTimeUrl,
		urlParameters,
	)

	if ouraError != nil {
		return SleepTimes{},
			ouraError
	}

	var sleepTimes SleepTimes
	err := json.Unmarshal(*apiResponse, &sleepTimes)
	if err != nil {
		return SleepTimes{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return sleepTimes, nil
}

func (c *Client) GetSleepTime(documentId string) (SleepTime, *OuraError) {

	apiResponse, ouraError := c.Getter(
		fmt.Sprintf(SleepTimeUrl+"/%s", documentId),
		nil,
	)

	if ouraError != nil {
		return SleepTime{},
			ouraError
	}

	var sleepTime SleepTime
	err := json.Unmarshal(*apiResponse, &sleepTime)
	if err != nil {
		return SleepTime{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return sleepTime, nil
}
