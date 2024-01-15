package go_oura

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type PersonalInfo struct {
	ID     string  `json:"id"`
	Age    int     `json:"age"`
	Height float32 `json:"height"`
	Weight float32 `json:"weight"`
	Sex    string  `json:"biological_sex"`
	Email  string  `json:"email"`
}

type personalInfoBase PersonalInfo

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

func (c *Client) GetPersonalInfo() (PersonalInfo, *OuraError) {

	apiResponse, ouraError := c.Getter(PersonalInfoUrl, nil)

	if ouraError != nil {
		return PersonalInfo{},
			ouraError
	}

	var personalInfo PersonalInfo
	err := json.Unmarshal(*apiResponse, &personalInfo)
	if err != nil {
		return PersonalInfo{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return personalInfo, nil
}
