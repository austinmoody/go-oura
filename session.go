package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type Sessions struct {
	Items      []Session `json:"data"`
	Next_Token string    `json:"next_token"`
}

type Session struct {
	ID                       string           `json:"id"`
	Day                      Date             `json:"day"`
	StartDatetime            time.Time        `json:"start_datetime"`
	EndDatetime              time.Time        `json:"end_datetime"`
	Type                     string           `json:"type"`
	HeartRateData            SessionDataItems `json:"heart_rate"`
	HeartRateVariabilityData SessionDataItems `json:"heart_rate_variability"`
	Mood                     string           `json:"mood"`
	MotionCountData          SessionDataItems `json:"motion_count"`
}

type SessionDataItems struct {
	Interval  float64   `json:"interval"`
	Items     []float64 `json:"items"`
	Timestamp time.Time `json:"timestamp"`
}

type sessionBase Session
type sessionsBase Sessions

func (s *Session) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*s), data); err != nil {
		return err
	}

	var session sessionBase
	err := json.Unmarshal(data, &session)
	if err != nil {
		return err
	}

	*s = Session(session)
	return nil
}

func (s *Sessions) UnmarshalJSON(data []byte) error {
	if err := checkJSONFields(reflect.TypeOf(*s), data); err != nil {
		return err
	}

	var sessions sessionsBase
	err := json.Unmarshal(data, &sessions)
	if err != nil {
		return err
	}

	*s = Sessions(sessions)

	return nil
}

func (c *Client) GetSession(documentId string) (Session, *OuraError) {

	apiResponse, ouraError := c.Getter(
		fmt.Sprintf(SessionUrl+"/%s", documentId),
		nil,
	)

	if ouraError != nil {
		return Session{},
			ouraError
	}

	var session Session
	err := json.Unmarshal(*apiResponse, &session)
	if err != nil {
		return Session{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return session, nil
}

func (c *Client) GetSessions(startDate time.Time, endDate time.Time, nextToken *string) (Sessions, *OuraError) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, ouraError := c.Getter(
		SessionUrl,
		urlParameters,
	)

	if ouraError != nil {
		return Sessions{},
			ouraError
	}

	var sessions Sessions
	err := json.Unmarshal(*apiResponse, &sessions)
	if err != nil {
		return Sessions{},
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to process response body with error: %v", err),
			}
	}

	return sessions, nil
}
