// Package go_oura provides a simple binding to the Oura Ring v2 API

// This file contains code related to sessions recorded by an Oura Ring
// Session API https://cloud.ouraring.com/v2/docs#tag/Session-Routes

package go_oura

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// Sessions stores a list of session items along with a token which may be used to pull the next batch of Session items from the API.
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Multiple_session_Documents_v2_usercollection_session_get
type Sessions struct {
	Items     []Session `json:"data"`
	NextToken string    `json:"next_token"`
}

// Session stores specifics for a single session recorded by an Oura Ring
// JSON described at https://cloud.ouraring.com/v2/docs#operation/Single_session_Document_v2_usercollection_session__document_id__get
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

type SessionDataItems IntervalItems

type sessionBase Session
type sessionsBase Sessions

// UnmarshalJSON is a helper function to convert a recorded session JSON from the API to the Session type.
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

// UnmarshalJSON is a helper function to convert multiple recorded sessions JSON from the API to the Sessions type.
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

// GetSession calls the Oura Ring API with a specific session identifier and returns a Session object
func (c *Client) GetSession(sessionId string) (Session, error) {

	apiResponse, err := c.Getter(
		fmt.Sprintf(SessionUrl+"/%s", sessionId),
		nil,
	)

	if err != nil {
		return Session{},
			err
	}

	var session Session
	err = json.Unmarshal(*apiResponse, &session)
	if err != nil {
		return Session{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return session, nil
}

// GetSessions accepts a start & end date and returns a Sessions object which will contain any Session
// found in the time period.  Optionally the next token can be passed which tells the API to give the next set of
// sessions if the date range returns a large set.
func (c *Client) GetSessions(startDate time.Time, endDate time.Time, nextToken *string) (Sessions, error) {

	urlParameters := url.Values{
		"start_date": []string{startDate.Format("2006-01-02")},
		"end_date":   []string{endDate.Format("2006-01-02")},
	}

	if nextToken != nil {
		urlParameters.Set("next_token", *nextToken)
	}

	apiResponse, err := c.Getter(
		SessionUrl,
		urlParameters,
	)

	if err != nil {
		return Sessions{},
			err
	}

	var sessions Sessions
	err = json.Unmarshal(*apiResponse, &sessions)
	if err != nil {
		return Sessions{}, fmt.Errorf("failed to process response body with error: %v", err)
	}

	return sessions, nil
}
