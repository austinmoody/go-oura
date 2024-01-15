package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type Activities struct {
	Items     []Activity `json:"data"`
	NextToken *string    `json:"next_token"`
}

type Activity struct {
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

type dailyActivityBase Activity
type dailyActivitiesBase Activities

func (da *Activities) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*da), data); err != nil {
		return err
	}

	var aBase dailyActivitiesBase
	err := json.Unmarshal(data, &aBase)
	if err != nil {
		return err
	}

	*da = Activities(aBase)
	return nil
}

func (da *Activity) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*da), data); err != nil {
		return err
	}

	var aBase dailyActivityBase
	err := json.Unmarshal(data, &aBase)
	if err != nil {
		return err
	}

	*da = Activity(aBase)
	return nil
}

func (c *Client) GetActivity(documentId string) (Activity, *OuraError) {
	apiResponse, ouraError := c.Getter(fmt.Sprintf(ActivityUrl+"/%s", documentId), nil)

	if ouraError != nil {
		return Activity{},
			ouraError
	}

	var activity Activity
	err := json.Unmarshal(*apiResponse, &activity)
	if err != nil {
		return Activity{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return activity, nil
}

func (c *Client) GetActivities(startDate time.Time, endDate time.Time, nextToken *string) (Activities, *OuraError) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		ActivityUrl,
		urlParameters,
	)

	if ouraError != nil {
		return Activities{}, ouraError
	}

	var activities Activities
	err := json.Unmarshal(*apiResponse, &activities)
	if err != nil {
		return Activities{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return activities, nil
}
