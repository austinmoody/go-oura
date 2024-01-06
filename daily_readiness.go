package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

/*
https://cloud.ouraring.com/v2/docs#tag/Daily-Readiness-Routes
*/

type DailyReadinessDocuments struct {
	Documents []DailyReadinessDocument `json:"data"`
	NextToken *string                  `json:"next_token"`
}

type DailyReadinessDocument struct {
	Id           string `json:"id"`
	Contributors struct {
		ActivityBalance     int `json:"activity_balance"`
		BodyTemperature     int `json:"body_temperature"`
		HrvBalance          int `json:"hrv_balance"`
		PreviousDayActivity int `json:"previous_day_activity"`
		PreviousNight       int `json:"previous_night"`
		RecoveryIndex       int `json:"recovery_index"`
		RestingHeartRate    int `json:"resting_heart_rate"`
		SleepBalance        int `json:"sleep_balance"`
	} `json:"contributors"`
	Day                       Date      `json:"day"`
	Score                     int       `json:"score"`
	TemperatureDeviation      float64   `json:"temperature_deviation"`
	TemperatureTrendDeviation float64   `json:"temperature_trend_deviation"`
	Timestamp                 time.Time `json:"timestamp"`
}

func (c *Client) GetReadinessDocuments(startDate time.Time, endDate time.Time) (DailyReadinessDocuments, error) {

	apiResponse, ouraError := c.Getter(
		"usercollection/daily_readiness",
		url.Values{
			"start_date": []string{startDate.Format("2006-01-02")},
			"end_date":   []string{endDate.Format("2006-01-02")},
		},
	)

	if ouraError != nil {
		return DailyReadinessDocuments{},
			fmt.Errorf("failed to get API response with error: %w", ouraError)
	}

	var readiness DailyReadinessDocuments
	err := json.Unmarshal(*apiResponse, &readiness)
	if err != nil {
		return DailyReadinessDocuments{},
			fmt.Errorf("failed to process response body with error: %w", err)
	}

	return readiness, nil
}

func (c *Client) GetReadinessDocument(documentId string) (DailyReadinessDocument, error) {

	apiResponse, ouraError := c.Getter(fmt.Sprintf("/usercollection/daily_readiness/%s", documentId), nil)

	if ouraError != nil {
		return DailyReadinessDocument{},
			fmt.Errorf("failed to get API response with error: %w", ouraError)
	}

	var readiness DailyReadinessDocument
	err := json.Unmarshal(*apiResponse, &readiness)
	if err != nil {
		return DailyReadinessDocument{},
			fmt.Errorf("failed to process response body with error: %w", err)
	}

	return readiness, nil
}
