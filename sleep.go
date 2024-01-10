package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type Sleeps struct {
	Items     []Sleep `json:"data"`
	NextToken string  `json:"next_token"`
}

type Sleep struct {
	ID           string            `json:"id"`
	Contributors SleepContributors `json:"contributors"`
	Day          Date              `json:"day"`
	Score        int64             `json:"score"`
	Timestamp    time.Time         `json:"timestamp"`
}

type SleepContributors struct {
	DeepSleep   int64 `json:"deep_sleep"`
	Efficiency  int64 `json:"efficiency"`
	Latency     int64 `json:"latency"`
	RemSleep    int64 `json:"rem_sleep"`
	Restfulness int64 `json:"restfulness"`
	Timing      int64 `json:"timing"`
	TotalSleep  int64 `json:"total_sleep"`
}

type dailySleepDocumentBase Sleep
type dailySleepDocumentsBase Sleeps

func (sd *Sleep) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*sd)
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

	var documentBase dailySleepDocumentBase
	err = json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*sd = Sleep(documentBase)
	return nil
}

func (sd *Sleeps) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*sd)
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

	var documentBase dailySleepDocumentsBase
	err = json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*sd = Sleeps(documentBase)
	return nil
}

func (c *Client) GetSleeps(startDate time.Time, endDate time.Time) (Sleeps, *OuraError) {
	apiResponse, ouraError := c.Getter(
		SleepUrl,
		url.Values{
			"start_date": []string{startDate.Format("2006-01-02")},
			"end_date":   []string{endDate.Format("2006-01-02")},
		},
	)

	if ouraError != nil {
		return Sleeps{},
			ouraError
	}

	var documents Sleeps
	err := json.Unmarshal(*apiResponse, &documents)
	if err != nil {
		return Sleeps{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return documents, nil
}

func (c *Client) GetSleep(documentId string) (Sleep, *OuraError) {

	apiResponse, ouraError := c.Getter(fmt.Sprintf(SleepUrl+"/%s", documentId), nil)

	if ouraError != nil {
		return Sleep{},
			ouraError
	}

	var sleep Sleep
	err := json.Unmarshal(*apiResponse, &sleep)
	if err != nil {
		return Sleep{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return sleep, nil
}
