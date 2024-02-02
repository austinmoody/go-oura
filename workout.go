// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to workouts recorded by an Oura Ring
// https://cloud.ouraring.com/v2/docs#tag/Workout-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// Workouts stores a list of workout items along with a token which may be used to pull the next batch of Workout items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_workout_Documents_v2_usercollection_workout_get
type Workouts struct {
	Items     []Workout `json:"data"`
	NextToken string    `json:"next_token"`
}

// Workout stores specifics for a single recorded workout
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_workout_Document_v2_usercollection_workout__document_id__get
type Workout struct {
	Id            string    `json:"id"`
	Activity      string    `json:"activity"`
	Calories      float64   `json:"calories"`
	Day           Date      `json:"day"`
	Distance      float64   `json:"distance"`
	EndDatetime   time.Time `json:"end_datetime"`
	Intensity     string    `json:"intensity"`
	Label         string    `json:"label"`
	Source        string    `json:"source"`
	StartDatetime time.Time `json:"start_datetime"`
}

type workoutBase Workout
type workoutsBase Workouts

// UnmarshalJSON is a helper function to convert a workout JSON from the API to the Workout type.
func (w *Workout) UnmarshalJSON(data []byte) error {

	if err := checkJSONFields(reflect.TypeOf(*w), data); err != nil {
		return err
	}

	var workout workoutBase
	err := json.Unmarshal(data, &workout)
	if err != nil {
		return err
	}

	*w = Workout(workout)

	return nil
}

// UnmarshalJSON is a helper function to convert multiple workouts JSON from the API to the Workouts type.
func (w *Workouts) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*w), data); err != nil {
		return err
	}

	var workouts workoutsBase
	err := json.Unmarshal(data, &workouts)
	if err != nil {
		return err
	}

	*w = Workouts(workouts)
	return nil
}

// GetWorkouts accepts a start & end date and returns a Workouts object which will contain any Workout
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// workouts if the date range returns a large set.
func (c *Client) GetWorkouts(startDate time.Time, endDate time.Time, nextToken *string) (Workouts, error) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, err := c.Getter(
		WorkoutUrl,
		urlParameters,
	)

	if err != nil {
		return Workouts{},
			err
	}

	var workouts Workouts
	err = json.Unmarshal(*apiResponse, &workouts)
	if err != nil {
		return Workouts{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return workouts, nil
}

// GetWorkout calls the Oura Ring API with a specific workout identifier and returns a Workout object
func (c *Client) GetWorkout(workoutId string) (Workout, error) {

	apiResponse, err := c.Getter(fmt.Sprintf(
		WorkoutUrl+"/%s",
		workoutId,
	), nil)

	if err != nil {
		return Workout{},
			err
	}

	var workout Workout
	err = json.Unmarshal(*apiResponse, &workout)
	if err != nil {
		return Workout{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return workout, nil
}
