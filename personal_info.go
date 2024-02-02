// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to a user's personal information
// Personal Info API description: https://cloud.ouraring.com/v2/docs#tag/Personal-Info-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// PersonalInfo stores the user's information
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_Personal_Info_Document_v2_usercollection_personal_info_get
type PersonalInfo struct {
	ID     string  `json:"id"`
	Age    int     `json:"age"`
	Height float32 `json:"height"`
	Weight float32 `json:"weight"`
	Sex    string  `json:"biological_sex"`
	Email  string  `json:"email"`
}

type personalInfoBase PersonalInfo

// UnmarshalJSON is a helper function to convert a personal info JSON from the API to the PersonalInfo type.
func (pi *PersonalInfo) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*pi), data); err != nil {
		return err
	}

	var personalInfo personalInfoBase
	err := json.Unmarshal(data, &personalInfo)
	if err != nil {
		return err
	}

	*pi = PersonalInfo(personalInfo)
	return nil

}

// GetPersonalInfo calls the Oura Ring API and returns a PersonalInfo object describing the user
func (c *Client) GetPersonalInfo() (PersonalInfo, error) {

	apiResponse, err := c.Getter(PersonalInfoUrl, nil)

	if err != nil {
		return PersonalInfo{},
			err
	}

	var personalInfo PersonalInfo
	err = json.Unmarshal(*apiResponse, &personalInfo)
	if err != nil {
		return PersonalInfo{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return personalInfo, nil
}
