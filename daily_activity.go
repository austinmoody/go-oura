package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type DailyActivities struct {
	Activities []DailyActivity `json:"data"`
	NextToken  *string         `json:"next_token"`
}

type DailyActivity struct {
	ID                        string      `json:"id"`
	Class5Min                 string      `json:"class_5_min"`
	Score                     int         `json:"score"`
	ActiveCalories            int         `json:"active_calories"`
	AverageMetMinutes         float64     `json:"average_met_minutes"`
	Contributors              Contributor `json:"contributors"`
	EquivalentWalkingDistance int         `json:"equivalent_walking_distance"`
	HighActivityMetMinutes    int         `json:"high_activity_met_minutes"`
	HighActivityTime          int         `json:"high_activity_time"`
	InactivityAlerts          int         `json:"inactivity_alerts"`
	LowActivityMetMinutes     int         `json:"low_activity_met_minutes"`
	LowActivityTime           int         `json:"low_activity_time"`
	MediumActivityMetMinutes  int         `json:"medium_activity_met_minutes"`
	MediumActivityTime        int         `json:"medium_activity_time"`
	Met                       Met         `json:"met"`
	MetersToTarget            int         `json:"meters_to_target"`
	NonWearTime               int         `json:"non_wear_time"`
	RestingTime               int         `json:"resting_time"`
	SedentaryMetMinutes       int         `json:"sedentary_met_minutes"`
	SedentaryTime             int         `json:"sedentary_time"`
	Steps                     int         `json:"steps"`
	TargetCalories            int         `json:"target_calories"`
	TargetMeters              int         `json:"target_meters"`
	TotalCalories             int         `json:"total_calories"`
	Day                       Date        `json:"day"`
	Timestamp                 time.Time   `json:"timestamp"`
}

type Contributor struct {
	MeetDailyTargets  int `json:"meet_daily_targets"`
	MoveEveryHour     int `json:"move_every_hour"`
	RecoveryTime      int `json:"recovery_time"`
	StayActive        int `json:"stay_active"`
	TrainingFrequency int `json:"training_frequency"`
	TrainingVolume    int `json:"training_volume"`
}

type Met struct {
	Interval  float64   `json:"interval"`
	Items     []float64 `json:"items"`
	Timestamp time.Time `json:"timestamp"`
}

type dailyActivityBase DailyActivity

func (da *DailyActivity) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*da)
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

	var aBase dailyActivityBase
	err = json.Unmarshal(data, &aBase)
	if err != nil {
		return err
	}

	*da = DailyActivity(aBase)
	return nil
}

func (c *Client) GetActivity(documentId string) (DailyActivity, error) {
	apiResponse, ouraError := c.Getter(fmt.Sprintf("/usercollection/daily_activity/%s", documentId), nil)

	if ouraError != nil {
		return DailyActivity{},
			fmt.Errorf("failed to get API response with error: %w", ouraError)
	}

	var activity DailyActivity
	err := json.Unmarshal(*apiResponse, &activity)
	if err != nil {
		return DailyActivity{},
			fmt.Errorf("failed to process response body with error: %w", err)
	}

	return activity, nil
}

func (c *Client) GetActivities(startDate time.Time, endDate time.Time) (DailyActivities, error) {

	apiResponse, ouraError := c.Getter(
		"usercollection/daily_activity",
		url.Values{
			"start_date": []string{startDate.Format("2006-01-02")},
			"end_date":   []string{endDate.Format("2006-01-02")},
		},
	)

	if ouraError != nil {
		return DailyActivities{}, fmt.Errorf("failed to get API response with error: %w", ouraError)
	}

	var activities DailyActivities
	err := json.Unmarshal(*apiResponse, &activities)
	if err != nil {
		return DailyActivities{}, fmt.Errorf("failed to process response body with error: %w", err)
	}

	return activities, nil
}
