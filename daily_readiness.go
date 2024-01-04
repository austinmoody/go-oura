package go_oura

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(input []byte) error {
	strInput := strings.Trim(string(input), `"`)
	newTime, err := time.Parse("2006-01-02", strInput)
	if err != nil {
		return err
	}

	d.Time = newTime
	return nil
}

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

	apiUrl, err := url.Parse(c.config.BaseUrl)
	if err != nil {
		return DailyReadinessDocuments{},
			fmt.Errorf("failed to parse base url with error: %w", err)
	}

	apiUrl.Path = path.Join(apiUrl.Path, "usercollection/daily_readiness")

	params := url.Values{}
	params.Add("start_date", startDate.Format("2006-01-02"))
	params.Add("end_date", endDate.Format("2006-01-02"))

	apiUrl.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return DailyReadinessDocuments{},
			fmt.Errorf("failed to create a new HTTP request with error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.config.accessToken)

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return DailyReadinessDocuments{},
			fmt.Errorf("failed to complete HTTP request with error: %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return DailyReadinessDocuments{},
			fmt.Errorf("failed to read response body with error: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return DailyReadinessDocuments{},
			fmt.Errorf("failed to close response body with error: %w", err)
	}

	var readiness DailyReadinessDocuments
	err = json.Unmarshal(data, &readiness)
	if err != nil {
		return DailyReadinessDocuments{},
			fmt.Errorf("failed to process response body with error: %w", err)
	}

	return readiness, nil

}

func (c *Client) GetReadinessDocument(documentId string) (DailyReadinessDocument, error) {
	apiUrl := fmt.Sprintf("%s/usercollection/daily_readiness/%s", c.config.BaseUrl, documentId)

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return DailyReadinessDocument{},
			fmt.Errorf("failed to create a new HTTP request with error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.config.accessToken)

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return DailyReadinessDocument{},
			fmt.Errorf("failed to complete HTTP request with error: %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return DailyReadinessDocument{},
			fmt.Errorf("failed to read response body with error: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return DailyReadinessDocument{},
			fmt.Errorf("failed to close response body with error: %w", err)
	}

	var readiness DailyReadinessDocument
	err = json.Unmarshal(data, &readiness)
	if err != nil {
		return DailyReadinessDocument{},
			fmt.Errorf("failed to process response body with error: %w", err)
	}

	return readiness, nil
}
