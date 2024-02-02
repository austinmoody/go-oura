// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to Enhanced Tags recorded via the Oura Ring Application
// Enhanced Tags API description: https://cloud.ouraring.com/v2/docs#tag/Enhanced-Tag-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// EnhancedTags stores a list of EnhancedTag items along with a token which may be used to pull the next batch of EnhancedTag items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_enhanced_tag_Documents_v2_usercollection_enhanced_tag_get
type EnhancedTags struct {
	Items     []EnhancedTag `json:"data"`
	NextToken string        `json:"next_token"`
}

// EnhancedTag describes a single tag value
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_enhanced_tag_Document_v2_usercollection_enhanced_tag__document_id__get
type EnhancedTag struct {
	ID          string     `json:"id"`
	TagTypeCode string     `json:"tag_type_code"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	StartDay    *Date      `json:"start_day"`
	EndDay      *Date      `json:"end_day"`
	Comment     string     `json:"comment"`
}

type tagDocumentBase EnhancedTag
type tagDocumentsBase EnhancedTags

// UnmarshalJSON is a helper function to convert enhanced tag JSON from the API to the EnhancedTag type.
func (t *EnhancedTag) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*t), data); err != nil {
		return err
	}

	var documentBase tagDocumentBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*t = EnhancedTag(documentBase)
	return nil
}

// UnmarshalJSON is a helper function to convert multiple enhanced tag JSON from the API to the EnhancedTags type.
func (t *EnhancedTags) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*t), data); err != nil {
		return err
	}

	var documentBase tagDocumentsBase
	err := json.Unmarshal(data, &documentBase)
	if err != nil {
		return err
	}

	*t = EnhancedTags(documentBase)
	return nil
}

// GetEnhancedTags accepts a start & end date and returns a EnhancedTags object which will contain any EnhancedTag items
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// tags if the date range returns a large set.
func (c *Client) GetEnhancedTags(startDate time.Time, endDate time.Time, nextToken *string) (EnhancedTags, error) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, err := c.Getter(
		TagUrl,
		urlParameters,
	)

	if err != nil {
		return EnhancedTags{},
			err
	}

	var documents EnhancedTags
	err = json.Unmarshal(*apiResponse, &documents)
	if err != nil {
		return EnhancedTags{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return documents, nil
}

// GetEnhancedTag accepts a single Enhanced Tag ID and returns a EnhancedTag object.
func (c *Client) GetEnhancedTag(documentId string) (EnhancedTag, error) {

	apiResponse, err := c.Getter(fmt.Sprintf(TagUrl+"/%s", documentId), nil)

	if err != nil {
		return EnhancedTag{},
			err
	}

	var tag EnhancedTag
	err = json.Unmarshal(*apiResponse, &tag)
	if err != nil {
		return EnhancedTag{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return tag, nil
}
