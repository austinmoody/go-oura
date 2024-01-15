package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type RingConfigurations struct {
	Items     []RingConfiguration `json:"data"`
	NextToken string              `json:"next_token"`
}

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

func (c *Client) GetRingConfigurations(startDate time.Time, endDate time.Time, nextToken *string) (RingConfigurations, *OuraError) {
	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		RingConfigurationUrl,
		urlParameters,
	)

	if ouraError != nil {
		return RingConfigurations{},
			ouraError
	}

	var ringConfigurations RingConfigurations
	err := json.Unmarshal(*apiResponse, &ringConfigurations)
	if err != nil {
		return RingConfigurations{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return ringConfigurations, nil
}

func (c *Client) GetRingConfiguration(documentId string) (RingConfiguration, *OuraError) {

	apiResponse, ouraError := c.Getter(
		fmt.Sprintf(RingConfigurationUrl+"/%s", documentId),
		nil,
	)

	if ouraError != nil {
		return RingConfiguration{},
			ouraError
	}

	var ringConfiguration RingConfiguration
	err := json.Unmarshal(*apiResponse, &ringConfiguration)
	if err != nil {
		return RingConfiguration{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return ringConfiguration, nil

}
