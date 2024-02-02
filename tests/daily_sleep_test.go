package tests

import (
	"github.com/austinmoody/go_oura"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetSleepDocument(t *testing.T) {
	tt := []struct {
		name           string
		documentId     string
		mockResponse   string
		expectedOutput go_oura.DailySleep
		expectErr      bool
	}{
		{
			name:         "Valid_SleepDocument_Response",
			documentId:   "1",
			mockResponse: `{"id":"4eaa0e18-3464-49cc-961a-2ffd5f8ea98e","contributors":{"deep_sleep":63,"efficiency":93,"latency":64,"rem_sleep":95,"restfulness":72,"timing":94,"total_sleep":90},"day":"2024-01-07","score":83,"timestamp":"2024-01-07T00:00:00+00:00"}`,
			expectedOutput: go_oura.DailySleep{
				ID: "4eaa0e18-3464-49cc-961a-2ffd5f8ea98e",
				Day: func() go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2024-01-07")
					return go_oura.Date{Time: t}
				}(),
				Score: 83,
				Timestamp: func() time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2024-01-07T00:00:00+00:00")
					return t
				}(),
				Contributors: go_oura.SleepContributors{
					DeepSleep:   63,
					Efficiency:  93,
					Latency:     64,
					RemSleep:    95,
					Restfulness: 72,
					Timing:      94,
					TotalSleep:  90,
				},
			},
			expectErr: false,
		},
		{
			name:           "Invalid_SleepDocument_Response",
			documentId:     "2",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.DailySleep{},
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

			activity, err := client.GetDailySleep(tc.documentId)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}

				return
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(activity, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, activity)
			}
		})
	}
}

func TestGetSleepDocuments(t *testing.T) {
	tt := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		mockResponse   string
		expectedOutput go_oura.DailySleeps
		expectErr      bool
	}{
		{
			name:         "Valid_SleepDocuments_Response",
			startTime:    time.Now().Add(-1 * time.Hour),
			endTime:      time.Now().Add(-2 * time.Hour),
			mockResponse: `{"data":[{"id":"4eaa0e18-3464-49cc-961a-2ffd5f8ea98e","contributors":{"deep_sleep":63,"efficiency":93,"latency":64,"rem_sleep":95,"restfulness":72,"timing":94,"total_sleep":90},"day":"2024-01-07","score":83,"timestamp":"2024-01-07T00:00:00+00:00"}],"next_token":null}`,
			expectedOutput: go_oura.DailySleeps{
				Items: []go_oura.DailySleep{
					{
						ID: "4eaa0e18-3464-49cc-961a-2ffd5f8ea98e",
						Day: func() go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2024-01-07")
							return go_oura.Date{Time: t}
						}(),
						Score: 83,
						Timestamp: func() time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2024-01-07T00:00:00+00:00")
							return t
						}(),
						Contributors: go_oura.SleepContributors{
							DeepSleep:   63,
							Efficiency:  93,
							Latency:     64,
							RemSleep:    95,
							Restfulness: 72,
							Timing:      94,
							TotalSleep:  90,
						},
					},
				},
			},
			expectErr: false,
		}, {
			name:           "Invalid_SleepDocuments_Response",
			startTime:      time.Now(),
			endTime:        time.Now(),
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.DailySleeps{},
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

			activity, err := client.GetDailySleeps(tc.startTime, tc.endTime, nil)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}

				return
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(activity, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, activity)
			}
		})
	}
}
