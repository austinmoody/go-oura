package go_oura

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type RestModes struct {
	Items     []RestMode `json:"data"`
	NextToken string     `json:"next_token"`
}

type RestMode struct {
	ID        string    `json:"id"`
	EndDay    Date      `json:"end_day"`
	EndTime   time.Time `json:"end_time"`
	Episodes  []Episode `json:"episodes"`
	StartDay  Date      `json:"start_day"`
	StartTime time.Time `json:"start_time"`
}

type Episode struct {
	Tags      []string  `json:"tags"`
	Timestamp time.Time `json:"timestamp"`
}

type restModeBase RestMode
type restModesBase RestModes

func (rm *RestMode) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*rm)
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

	var restMode restModeBase
	err = json.Unmarshal(data, &restMode)
	if err != nil {
		return err
	}

	*rm = RestMode(restMode)
	return nil
}

func (rm *RestModes) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*rm)
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

	var restModes restModesBase
	err = json.Unmarshal(data, &restModes)
	if err != nil {
		return err
	}

	*rm = RestModes(restModes)
	return nil
}

func (c *Client) GetRestMode(documentId string) (RestMode, *OuraError) {

	apiResponse, ouraError := c.Getter(fmt.Sprintf(RestModeUrl+"/%s", documentId), nil)

	if ouraError != nil {
		return RestMode{},
			ouraError
	}

	var restMode RestMode
	err := json.Unmarshal(*apiResponse, &restMode)
	if err != nil {
		return RestMode{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return restMode, nil
}
