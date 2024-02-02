// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to ring configurations, basically information about your rings
// Ring Configuration API https://cloud.ouraring.com/v2/docs#tag/Ring-Configuration-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// RingConfigurations stores a list of ring configuration items along with a token which may be used to pull the next batch of RingConfiguration items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_ring_configuration_Documents_v2_usercollection_ring_configuration_get
type RingConfigurations struct {
	Items     []RingConfiguration `json:"data"`
	NextToken string              `json:"next_token"`
}

// RingConfiguration stores specifics for a single Ring's configuration
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_ring_configuration_Document_v2_usercollection_ring_configuration__document_id__get
type RingConfiguration struct {
	ID              string     `json:"id"`
	Color           string     `json:"color"`
	Design          string     `json:"design"`
	FirmwareVersion string     `json:"firmware_version"`
	HardwareType    string     `json:"hardware_type"`
	SetUpAt         *time.Time `json:"set_up_at"`
	Size            int        `json:"size"`
}

type ringConfigurationBase RingConfiguration
type ringConfigurationsBase RingConfigurations

// UnmarshalJSON is a helper function to convert a ring configuration JSON from the API to the RingConfiguration type.
func (rc *RingConfiguration) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*rc), data); err != nil {
		return err
	}

	var ringConfiguration ringConfigurationBase
	err := json.Unmarshal(data, &ringConfiguration)
	if err != nil {
		return err
	}

	*rc = RingConfiguration(ringConfiguration)
	return nil
}

// UnmarshalJSON is a helper function to convert a list of ring configurations JSON from the API to the RingConfigurations type.
func (rc *RingConfigurations) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*rc), data); err != nil {
		return err
	}

	var ringConfigurations ringConfigurationsBase
	err := json.Unmarshal(data, &ringConfigurations)
	if err != nil {
		return err
	}

	*rc = RingConfigurations(ringConfigurations)
	return nil
}

// GetRingConfigurations accepts a start & end date and returns a RingConfigurations object which will contain any RingConfiguration
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// ring configurations if the date range returns a large set.
func (c *Client) GetRingConfigurations(startDate time.Time, endDate time.Time, nextToken *string) (RingConfigurations, error) {
	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, err := c.Getter(
		RingConfigurationUrl,
		urlParameters,
	)

	if err != nil {
		return RingConfigurations{},
			err
	}

	var ringConfigurations RingConfigurations
	err = json.Unmarshal(*apiResponse, &ringConfigurations)
	if err != nil {
		return RingConfigurations{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return ringConfigurations, nil
}

// GetRingConfiguration calls the Oura Ring API with a specific Ring Configuration identifier and returns a RingConfiguration object
func (c *Client) GetRingConfiguration(ringConfigurationId string) (RingConfiguration, error) {

	apiResponse, err := c.Getter(
		fmt.Sprintf(RingConfigurationUrl+"/%s", ringConfigurationId),
		nil,
	)

	if err != nil {
		return RingConfiguration{},
			err
	}

	var ringConfiguration RingConfiguration
	err = json.Unmarshal(*apiResponse, &ringConfiguration)
	if err != nil {
		return RingConfiguration{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return ringConfiguration, nil

}
