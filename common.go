// This file contains any data types which are used across different Oura Ring types.

package go_oura

import (
	"encoding/json"
	"reflect"
	"time"
)

// IntervalItems is a common data type used by different Oura Ring types.  It stores an interval, a timestamp, and
// then 0 or more items recorded during that interval.
type IntervalItems struct {
	Interval  float64   `json:"interval"`
	Items     []float64 `json:"items"`
	Timestamp time.Time `json:"timestamp"`
}

type intervalItemsBase IntervalItems

func (ii *IntervalItems) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*ii), data); err != nil {
		return err
	}

	var documentBase intervalItemsBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*ii = IntervalItems(documentBase)
	return nil
}

// Contributors is a common data type used by different Oura Ring types and describes data points which contribute
// to a score.
type Contributors struct {
	ActivityBalance     int `json:"activity_balance"`
	BodyTemperature     int `json:"body_temperature"`
	HrvBalance          int `json:"hrv_balance"`
	PreviousDayActivity int `json:"previous_day_activity"`
	PreviousNight       int `json:"previous_night"`
	RecoveryIndex       int `json:"recovery_index"`
	RestingHeartRate    int `json:"resting_heart_rate"`
	SleepBalance        int `json:"sleep_balance"`
}

type contributorsBase Contributors

func (c *Contributors) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*c), data); err != nil {
		return err
	}

	var documentBase contributorsBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*c = Contributors(documentBase)
	return nil
}
