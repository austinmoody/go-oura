package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type Tags struct {
	Items     []Tag  `json:"data"`
	NextToken string `json:"next_token"`
}

type Tag struct {
	ID          string    `json:"id"`
	TagTypeCode string    `json:"tag_type_code"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	StartDay    string    `json:"start_day"`
	EndDay      string    `json:"end_day"`
	Comment     string    `json:"comment"`
}

type tagDocumentBase Tag
type tagDocumentsBase Tags

func (t *Tag) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	tt := reflect.TypeOf(*t)
	requiredFields := make([]string, 0, tt.NumField())
	for i := 0; i < tt.NumField(); i++ {
		jsonTag := tt.Field(i).Tag.Get("json")
		requiredFields = append(requiredFields, jsonTag)
	}

	for _, field := range requiredFields {
		if _, ok := rawMap[field]; !ok {
			return fmt.Errorf("required field %s not found", field)
		}
	}

	var documentBase tagDocumentBase
	err = json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*t = Tag(documentBase)
	return nil
}

func (t *Tags) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	tt := reflect.TypeOf(*t)
	requiredFields := make([]string, 0, tt.NumField())
	for i := 0; i < tt.NumField(); i++ {
		jsonTag := tt.Field(i).Tag.Get("json")
		requiredFields = append(requiredFields, jsonTag)
	}

	for _, field := range requiredFields {
		if _, ok := rawMap[field]; !ok {
			return fmt.Errorf("required field %s not found", field)
		}
	}

	var documentBase tagDocumentsBase
	err = json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*t = Tags(documentBase)
	return nil
}

func (c *Client) GetTags(startDate time.Time, endDate time.Time, nextToken *string) (Tags, *OuraError) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		TagUrl,
		urlParameters,
	)

	if ouraError != nil {
		return Tags{},
			ouraError
	}

	var documents Tags
	err := json.Unmarshal(*apiResponse, &documents)
	if err != nil {
		return Tags{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return documents, nil
}

func (c *Client) GetTag(documentId string) (Tag, *OuraError) {

	apiResponse, ouraError := c.Getter(fmt.Sprintf(TagUrl+"/%s", documentId), nil)

	if ouraError != nil {
		return Tag{},
			ouraError
	}

	var tag Tag
	err := json.Unmarshal(*apiResponse, &tag)
	if err != nil {
		return Tag{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return tag, nil
}
