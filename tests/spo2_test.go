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

func TestGetSpo2Reading(t *testing.T) {
	tt := []struct {
		name           string
		documentId     string
		mockResponse   string
		expectedOutput go_oura.Spo2Reading
		expectErr      bool
	}{
		{
			name:         "Valid_Spo2_Response",
			documentId:   "1",
			mockResponse: `{"id":"324f5ba1-f3f7-410a-b41e-c6585d1eaacc","day":"2024-01-09","spo2_percentage":{"average":98.781}}`,
			expectedOutput: go_oura.Spo2Reading{
				ID: "324f5ba1-f3f7-410a-b41e-c6585d1eaacc",
				Day: func() go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2024-01-09")
					return go_oura.Date{Time: t}
				}(),
				Percentage: go_oura.Spo2Percentage{
					Average: 98.781,
				},
			},
			expectErr: false,
		},
		{
			name:           "Invalid_Spo2_Response",
			documentId:     "2",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Spo2Reading{},
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

			activity, err := client.GetSpo2Reading(tc.documentId)
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

func TestGetSpo2Readings(t *testing.T) {
	tt := []struct {
		name           string
		startDate      time.Time
		endDate        time.Time
		mockResponse   string
		expectedOutput go_oura.Spo2Readings
		expectErr      bool
	}{
		{
			name:         "Valid_Multiple_Spo2_Response",
			startDate:    time.Now().Add(-1 * time.Hour),
			endDate:      time.Now().Add(-2 * time.Hour),
			mockResponse: `{"data":[{"id":"324f5ba1-f3f7-410a-b41e-c6585d1eaacc","day":"2024-01-09","spo2_percentage":{"average":98.781}},{"id":"4c2509bd-4a67-4694-9a62-e0582c2c34a2","day":"2024-01-10","spo2_percentage":{"average":98.688}}],"next_token":null}`,
			expectedOutput: go_oura.Spo2Readings{
				Items: []go_oura.Spo2Reading{
					{
						ID: "324f5ba1-f3f7-410a-b41e-c6585d1eaacc",
						Day: func() go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2024-01-09")
							return go_oura.Date{Time: t}
						}(),
						Percentage: go_oura.Spo2Percentage{
							Average: 98.781,
						},
					},
					{
						ID: "4c2509bd-4a67-4694-9a62-e0582c2c34a2",
						Day: func() go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2024-01-10")
							return go_oura.Date{Time: t}
						}(),
						Percentage: go_oura.Spo2Percentage{
							Average: 98.688,
						},
					},
				},
			},
			expectErr: false,
		}, {
			name:           "Invalid_Multiple_Spo2_Response",
			startDate:      time.Now().Add(-3 * time.Hour),
			endDate:        time.Now().Add(-4 * time.Hour),
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Spo2Readings{},
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

			activity, err := client.GetSpo2Readings(tc.startDate, tc.endDate)
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
