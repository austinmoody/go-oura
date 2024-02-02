// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to Daily Activities recorded by the Oura Ring
// Daily Activities API description: https://cloud.ouraring.com/v2/docs#tag/Daily-Activity-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// DailyActivities stores a list of daily activity items along with a token which may be used to pull the next batch of DailyActivity items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_daily_activity_Documents_v2_usercollection_daily_activity_get
type DailyActivities struct {
	Items     []DailyActivity `json:"data"`
	NextToken *string         `json:"next_token"`
}

// DailyActivity describes daily activity summary values and detailed activity levels.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_daily_activity_Document_v2_usercollection_daily_activity__document_id__get
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

// Contributor describes data points which contribute to the summary DailyActivity score
type Contributor struct {
	MeetDailyTargets  int `json:"meet_daily_targets"`
	MoveEveryHour     int `json:"move_every_hour"`
	RecoveryTime      int `json:"recovery_time"`
	StayActive        int `json:"stay_active"`
	TrainingFrequency int `json:"training_frequency"`
	TrainingVolume    int `json:"training_volume"`
}

// Met is a Metabolic Equivalent of Task Minutes.
type Met IntervalItems

type dailyActivityBase DailyActivity
type dailyActivitiesBase DailyActivities

// UnmarshalJSON is a helper function to convert daily activities JSON from the API to the DailyActivities type.
func (da *DailyActivities) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*da), data); err != nil {
		return err
	}

	var aBase dailyActivitiesBase
	err := json.Unmarshal(data, &aBase)
	if err != nil {
		return err
	}

	*da = DailyActivities(aBase)
	return nil
}

// UnmarshalJSON is a helper function to convert a daily activity JSON from the API to the DailyActivity type.
func (da *DailyActivity) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*da), data); err != nil {
		return err
	}

	var aBase dailyActivityBase
	err := json.Unmarshal(data, &aBase)
	if err != nil {
		return err
	}

	*da = DailyActivity(aBase)
	return nil
}

// GetActivity accepts a single Daily Activity ID and returns a DailyActivity object.
func (c *Client) GetActivity(dailyActivityId string) (DailyActivity, error) {
	apiResponse, err := c.Getter(
		fmt.Sprintf(ActivityUrl+"/%s", dailyActivityId),
		nil,
	)

	if err != nil {
		return DailyActivity{},
			err
	}

	var activity DailyActivity
	err = json.Unmarshal(*apiResponse, &activity)
	if err != nil {
		return DailyActivity{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return activity, nil
}

// GetActivities accepts a start & end date and returns a DailyActivities object which will contain any DailyActivity
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// activities if the date range returns a large set.
func (c *Client) GetActivities(startDate time.Time, endDate time.Time, nextToken *string) (DailyActivities, error) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, err := c.Getter(
		ActivityUrl,
		urlParameters,
	)

	if err != nil {
		return DailyActivities{}, err
	}

	var activities DailyActivities
	err = json.Unmarshal(*apiResponse, &activities)
	if err != nil {
		return DailyActivities{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return activities, nil
}
