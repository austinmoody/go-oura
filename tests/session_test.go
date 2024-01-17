package tests

import (
	"errors"
	go_oura "github.com/austinmoody/go_oura"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetSessionDocument(t *testing.T) {
	tt := []struct {
		name           string
		documentId     string
		mockResponse   string
		expectedOutput go_oura.Session
		expectErr      bool
	}{
		{
			name:         "Valid_Session_Response",
			documentId:   "valid-document",
			mockResponse: `{"id":"b630050a-f3d6-4feb-a04e-f1df2866fbf0","day":"2023-03-13","start_datetime":"2023-03-13T16:25:11-04:00","end_datetime":"2023-03-13T16:31:22-04:00","type":"meditation","heart_rate":{"interval":5,"items":[null,57.8,58],"timestamp":"2023-03-13T16:25:11.000-04:00"},"heart_rate_variability":{"interval":5,"items":[null,31,31.5],"timestamp":"2023-03-13T16:25:11.000-04:00"},"mood":"good","motion_count":{"interval":5,"items":[0,27,null],"timestamp":"2023-03-13T16:25:11.000-04:00"}}`,
			expectedOutput: go_oura.Session{
				ID: "b630050a-f3d6-4feb-a04e-f1df2866fbf0",
				Day: func() go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2023-03-13")
					return go_oura.Date{Time: t}
				}(),
				StartDatetime: func() time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2023-03-13T16:25:11-04:00")
					return t
				}(),
				EndDatetime: func() time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2023-03-13T16:31:22-04:00")
					return t
				}(),
				Type: "meditation",
				HeartRateData: go_oura.SessionDataItems{
					Interval: 5,
					Items:    []float64{0, 57.8, 58},
					Timestamp: func() time.Time {
						layout := "2006-01-02T15:04:05Z07:00"
						t, _ := time.Parse(layout, "2023-03-13T16:25:11.000-04:00")
						return t
					}(),
				},
				HeartRateVariabilityData: go_oura.SessionDataItems{
					Interval: 5,
					Items:    []float64{0, 31, 31.5},
					Timestamp: func() time.Time {
						layout := "2006-01-02T15:04:05Z07:00"
						t, _ := time.Parse(layout, "2023-03-13T16:25:11.000-04:00")
						return t
					}(),
				},
				Mood: "good",
				MotionCountData: go_oura.SessionDataItems{
					Interval: 5,
					Items:    []float64{0, 27, 0},
					Timestamp: func() time.Time {
						layout := "2006-01-02T15:04:05Z07:00"
						t, _ := time.Parse(layout, "2023-03-13T16:25:11.000-04:00")
						return t
					}(),
				},
			},
			expectErr: false,
		},
		{
			name:           "Invalid_Session_Response",
			documentId:     "invalid-session",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Session{},
			expectErr:      true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				_, err := rw.Write([]byte(tc.mockResponse))
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
			}))

			client := go_oura.NewClientWithUrlAndHttp("", server.URL, server.Client())

			session, err := client.GetSession(tc.documentId)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}

				var ouraErr *go_oura.OuraError
				if !errors.As(err, &ouraErr) {
					t.Errorf("expected an OuraError but got a different error: %v", err)
				}

				return
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(session, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, session)
			}
		})
	}
}

func TestGetSessions(t *testing.T) {
	tt := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		mockResponse   string
		expectedOutput go_oura.Sessions
		expectErr      bool
	}{
		{
			name:         "Valid_Sessions_Response",
			startTime:    time.Now().Add(-1 * time.Hour),
			endTime:      time.Now().Add(-2 * time.Hour),
			mockResponse: `{"data":[{"id":"b630050a-f3d6-4feb-a04e-f1df2866fbf0","day":"2023-03-13","start_datetime":"2023-03-13T16:25:11-04:00","end_datetime":"2023-03-13T16:31:22-04:00","type":"meditation","heart_rate":{"interval":5,"items":[null,57.8,58],"timestamp":"2023-03-13T16:25:11.000-04:00"},"heart_rate_variability":{"interval":5,"items":[null,31,31.5],"timestamp":"2023-03-13T16:25:11.000-04:00"},"mood":"good","motion_count":{"interval":5,"items":[0,27,null],"timestamp":"2023-03-13T16:25:11.000-04:00"}}],"next_token":"here-is-your-token"}`,
			expectedOutput: go_oura.Sessions{
				Items: []go_oura.Session{
					{
						ID: "b630050a-f3d6-4feb-a04e-f1df2866fbf0",
						Day: func() go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2023-03-13")
							return go_oura.Date{Time: t}
						}(),
						StartDatetime: func() time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2023-03-13T16:25:11-04:00")
							return t
						}(),
						EndDatetime: func() time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2023-03-13T16:31:22-04:00")
							return t
						}(),
						Type: "meditation",
						HeartRateData: go_oura.SessionDataItems{
							Interval: 5,
							Items:    []float64{0, 57.8, 58},
							Timestamp: func() time.Time {
								layout := "2006-01-02T15:04:05Z07:00"
								t, _ := time.Parse(layout, "2023-03-13T16:25:11.000-04:00")
								return t
							}(),
						},
						HeartRateVariabilityData: go_oura.SessionDataItems{
							Interval: 5,
							Items:    []float64{0, 31, 31.5},
							Timestamp: func() time.Time {
								layout := "2006-01-02T15:04:05Z07:00"
								t, _ := time.Parse(layout, "2023-03-13T16:25:11.000-04:00")
								return t
							}(),
						},
						Mood: "good",
						MotionCountData: go_oura.SessionDataItems{
							Interval: 5,
							Items:    []float64{0, 27, 0},
							Timestamp: func() time.Time {
								layout := "2006-01-02T15:04:05Z07:00"
								t, _ := time.Parse(layout, "2023-03-13T16:25:11.000-04:00")
								return t
							}(),
						},
					},
				},
				Next_Token: "here-is-your-token",
			},
			expectErr: false,
		}, {
			name:           "Invalid_Sessions_Response",
			startTime:      time.Now().Add(-3 * time.Hour),
			endTime:        time.Now().Add(-4 * time.Hour),
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Sessions{},
			expectErr:      true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				_, err := rw.Write([]byte(tc.mockResponse))
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
			}))

			client := go_oura.NewClientWithUrlAndHttp("", server.URL, server.Client())

			sessions, err := client.GetSessions(tc.startTime, tc.endTime, nil)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}

				var ouraErr *go_oura.OuraError
				if !errors.As(err, &ouraErr) {
					t.Errorf("expected an OuraError but got a different error: %v", err)
				}

				return
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(sessions, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, sessions)
			}
		})
	}
}
