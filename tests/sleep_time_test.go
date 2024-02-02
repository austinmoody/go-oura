package tests

import (
	"github.com/austinmoody/go_oura"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetSleepTime(t *testing.T) {
	tt := []struct {
		name           string
		documentId     string
		mockResponse   string
		expectedOutput go_oura.SleepTime
		expectErr      bool
	}{
		{
			name:         "Valid_SleepTime_Response_Without_Optimal",
			documentId:   "1",
			mockResponse: `{"id":"bb1044c6-6d85-406b-9bcd-0ce7dd438608","day":"2024-01-12","optimal_bedtime":null,"recommendation":"earlier_bedtime","status":"only_recommended_found"}`,
			expectedOutput: go_oura.SleepTime{
				ID: "bb1044c6-6d85-406b-9bcd-0ce7dd438608",
				Day: func() go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2024-01-12")
					return go_oura.Date{Time: t}
				}(),
				OptimalBedtime: nil,
				Recommendation: "earlier_bedtime",
				Status:         "only_recommended_found",
			},
		},
		{
			name:         "Valid_SleepTime_Response_With_Optimal",
			documentId:   "valid-with-optimal",
			mockResponse: `{"id":"bb1044c6-6d85-406b-9bcd-0ce7dd438608","day":"2024-01-12","optimal_bedtime":{"day_tz":1,"end_offset":2,"start_offset":3},"recommendation":"earlier_bedtime","status":"only_recommended_found"}`,
			expectedOutput: go_oura.SleepTime{
				ID: "bb1044c6-6d85-406b-9bcd-0ce7dd438608",
				Day: func() go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2024-01-12")
					return go_oura.Date{Time: t}
				}(),
				OptimalBedtime: &go_oura.OptimalBedtime{
					DayTz:       1,
					EndOffset:   2,
					StartOffset: 3,
				},
				Recommendation: "earlier_bedtime",
				Status:         "only_recommended_found",
			},
			expectErr: false,
		},
		{
			name:           "Invalid_SleepTime_Response",
			documentId:     "2",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.SleepTime{},
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

			activity, err := client.GetSleepTime(tc.documentId)
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

func TestGetSleepTimes(t *testing.T) {
	tt := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		mockResponse   string
		expectedOutput go_oura.SleepTimes
		expectErr      bool
	}{
		{
			name:         "Valid_SleepTimes_Response",
			startTime:    time.Now().Add(-1 * time.Hour),
			endTime:      time.Now().Add(-2 * time.Hour),
			mockResponse: `{"data":[{"id":"bb1044c6-6d85-406b-9bcd-0ce7dd438608","day":"2024-01-12","optimal_bedtime":null,"recommendation":"earlier_bedtime","status":"only_recommended_found"}],"next_token":"this-is-the-next-token"}`,
			expectedOutput: go_oura.SleepTimes{
				Items: []go_oura.SleepTime{
					{
						ID: "bb1044c6-6d85-406b-9bcd-0ce7dd438608",
						Day: func() go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2024-01-12")
							return go_oura.Date{Time: t}
						}(),
						OptimalBedtime: nil,
						Recommendation: "earlier_bedtime",
						Status:         "only_recommended_found",
					},
				},
				NextToken: "this-is-the-next-token",
			},
			expectErr: false,
		}, {
			name:           "Invalid_SleepTimes_Response",
			startTime:      time.Now().Add(-3 * time.Hour),
			endTime:        time.Now().Add(-4 * time.Hour),
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.SleepTimes{},
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

			activity, err := client.GetSleepTimes(tc.startTime, tc.endTime, nil)
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
