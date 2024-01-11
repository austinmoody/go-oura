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

func TestGetReadinesses(t *testing.T) {
	tt := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		mockResponse   string
		expectedOutput go_oura.Readinesses
		expectErr      bool
	}{
		{
			name:         "Valid_Readinesses_Response",
			startTime:    time.Now().Add(-1 * time.Hour),
			endTime:      time.Now().Add(-2 * time.Hour),
			mockResponse: `{"data":[{"id":"29a809a2-778c-4742-b945-e01876b8f32a","contributors":{"activity_balance":86,"body_temperature":89,"hrv_balance":4,"previous_day_activity":88,"previous_night":87,"recovery_index":99,"resting_heart_rate":1,"sleep_balance":90},"day":"2024-01-01","score":63,"temperature_deviation":0.2,"temperature_trend_deviation":0.38,"timestamp":"2024-01-01T00:00:00+00:00"}],"next_token":null}`,
			expectedOutput: go_oura.Readinesses{
				Items: []go_oura.Readiness{
					{
						Id: "29a809a2-778c-4742-b945-e01876b8f32a",
						Day: func() go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2024-01-01")
							return go_oura.Date{Time: t}
						}(),
						Score:                     63,
						TemperatureDeviation:      0.2,
						TemperatureTrendDeviation: 0.38,
						Timestamp: func() time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2024-01-01T00:00:00+00:00")
							return t
						}(),

						Contributors: go_oura.Contributors{
							ActivityBalance:     86,
							BodyTemperature:     89,
							HrvBalance:          4,
							PreviousDayActivity: 88,
							PreviousNight:       87,
							RecoveryIndex:       99,
							RestingHeartRate:    1,
							SleepBalance:        90,
						}},
				},
			},
			expectErr: false,
		},
		{
			name:           "Invalid_Readinesses_Response",
			startTime:      time.Now().Add(-3 * time.Hour),
			endTime:        time.Now().Add(-4 * time.Hour),
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Readinesses{},
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

			activity, err := client.GetReadinesses(tc.startTime, tc.endTime)
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

			if !reflect.DeepEqual(activity, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, activity)
			}
		})
	}

}

func TestGetReadiness(t *testing.T) {
	tt := []struct {
		name           string
		documentId     string
		mockResponse   string
		expectedOutput go_oura.Readiness
		expectErr      bool
	}{
		{
			name:         "Valid_Readiness_Response",
			documentId:   "1",
			mockResponse: `{"id":"29a809a2-778c-4742-b945-e01876b8f32a","contributors":{"activity_balance":86,"body_temperature":89,"hrv_balance":4,"previous_day_activity":88,"previous_night":87,"recovery_index":99,"resting_heart_rate":1,"sleep_balance":90},"day":"2024-01-01","score":63,"temperature_deviation":0.2,"temperature_trend_deviation":0.38,"timestamp":"2024-01-01T00:00:00+00:00"}`,
			expectedOutput: go_oura.Readiness{
				Id: "29a809a2-778c-4742-b945-e01876b8f32a",
				Day: func() go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2024-01-01")
					return go_oura.Date{Time: t}
				}(),
				Score:                     63,
				TemperatureDeviation:      0.2,
				TemperatureTrendDeviation: 0.38,
				Timestamp: func() time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2024-01-01T00:00:00+00:00")
					return t
				}(),

				Contributors: go_oura.Contributors{
					ActivityBalance:     86,
					BodyTemperature:     89,
					HrvBalance:          4,
					PreviousDayActivity: 88,
					PreviousNight:       87,
					RecoveryIndex:       99,
					RestingHeartRate:    1,
					SleepBalance:        90,
				},
			},
			expectErr: false,
		},
		{
			name:           "Invalid_Readiness_Response",
			documentId:     "2",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Readiness{},
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

			activity, err := client.GetReadiness(tc.documentId)
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

			if !reflect.DeepEqual(activity, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, activity)
			}
		})
	}
}
