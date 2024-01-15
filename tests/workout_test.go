package tests

import (
	"errors"
	"github.com/austinmoody/go-oura"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetWorkout(t *testing.T) {
	tt := []struct {
		name           string
		documentId     string
		mockResponse   string
		expectedOutput go_oura.Workout
		expectErr      bool
	}{
		{
			name:         "Valid_WorkoutDocument_Response",
			documentId:   "1",
			mockResponse: `{"id":"f07231a9-3600-40b1-aca7-9768e4284a5d","activity":"houseWork","calories":29.627,"day":"2024-01-06","distance":46.89835310974273,"end_datetime":"2024-01-06T09:27:00-05:00","intensity":"moderate","label":null,"source":"confirmed","start_datetime":"2024-01-06T09:14:00-05:00"}`,
			expectedOutput: go_oura.Workout{
				Id:       "f07231a9-3600-40b1-aca7-9768e4284a5d",
				Activity: "houseWork",
				Calories: 29.627,
				Day: func() go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2024-01-06")
					return go_oura.Date{Time: t}
				}(),
				Distance: 46.89835310974273,
				EndDatetime: func() time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2024-01-06T09:27:00-05:00")
					return t
				}(),
				Intensity: "moderate",
				Label:     "",
				Source:    "confirmed",
				StartDatetime: func() time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2024-01-06T09:14:00-05:00")
					return t
				}(),
			},
			expectErr: false,
		},
		{
			name:           "Invalid_WorkoutDocument_Response",
			documentId:     "2",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Workout{},
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

			workout, err := client.GetWorkout(tc.documentId)
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

			if !reflect.DeepEqual(workout, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, workout)
			}
		})
	}
}

func TestGetWorkouts(t *testing.T) {
	tt := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		mockResponse   string
		expectedOutput go_oura.Workouts
		expectErr      bool
	}{
		{
			name:         "Valid_Workouts_Response",
			startTime:    time.Now().Add(-1 * time.Hour),
			endTime:      time.Now().Add(-2 * time.Hour),
			mockResponse: `{"data":[{"id":"f07231a9-3600-40b1-aca7-9768e4284a5d","activity":"houseWork","calories":29.627,"day":"2024-01-06","distance":46.89835310974273,"end_datetime":"2024-01-06T09:27:00-05:00","intensity":"moderate","label":"your-label","source":"confirmed","start_datetime":"2024-01-06T09:14:00-05:00"}],"next_token":"your-token"}`,
			expectedOutput: go_oura.Workouts{
				Items: []go_oura.Workout{
					{
						Id:       "f07231a9-3600-40b1-aca7-9768e4284a5d",
						Activity: "houseWork",
						Calories: 29.627,
						Day: func() go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2024-01-06")
							return go_oura.Date{Time: t}
						}(),
						Distance: 46.89835310974273,
						EndDatetime: func() time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2024-01-06T09:27:00-05:00")
							return t
						}(),
						Intensity: "moderate",
						Label:     "your-label",
						Source:    "confirmed",
						StartDatetime: func() time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2024-01-06T09:14:00-05:00")
							return t
						}(),
					},
				},
				NextToken: "your-token",
			},
			expectErr: false,
		},
		{
			name:           "Invalid_Workouts_Response",
			startTime:      time.Now().Add(-3 * time.Hour),
			endTime:        time.Now().Add(-4 * time.Hour),
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Workouts{},
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

			workouts, err := client.GetWorkouts(tc.startTime, tc.endTime, nil)
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

			if !reflect.DeepEqual(workouts, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, workouts)
			}
		})
	}
}
