package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type Workouts struct {
	Items     []Workout `json:"data"`
	NextToken string    `json:"next_token"`
}

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

func (c *Client) GetWorkouts(startDate time.Time, endDate time.Time, nextToken *string) (Workouts, *OuraError) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		WorkoutUrl,
		urlParameters,
	)

	if ouraError != nil {
		return Workouts{},
			ouraError
	}

	var workouts Workouts
	err := json.Unmarshal(*apiResponse, &workouts)
	if err != nil {
		return Workouts{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return workouts, nil
}

func (c *Client) GetWorkout(documentId string) (Workout, *OuraError) {

	apiResponse, ouraError := c.Getter(fmt.Sprintf(
		WorkoutUrl+"/%s",
		documentId,
	), nil)

	if ouraError != nil {
		return Workout{},
			ouraError
	}

	var workout Workout
	err := json.Unmarshal(*apiResponse, &workout)
	if err != nil {
		return Workout{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return workout, nil
}
