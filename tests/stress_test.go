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

func TestGetStressDocument(t *testing.T) {
	tt := []struct {
		name           string
		documentId     string
		mockResponse   string
		expectedOutput go_oura.Stress
		expectErr      bool
	}{
		{
			name:         "Valid_Stress_Response",
			documentId:   "1",
			mockResponse: `{"id":"c9e6a9a9-2af3-4284-bbb7-038346c06bc9","day":"2024-01-11","stress_high":6300,"recovery_high":1800,"day_summary":"normal"}`,
			expectedOutput: go_oura.Stress{
				ID: "c9e6a9a9-2af3-4284-bbb7-038346c06bc9",
				Day: func() go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2024-01-11")
					return go_oura.Date{Time: t}
				}(),
				StressHigh:   6300,
				RecoveryHigh: 1800,
				DaySummary:   "normal",
			},
			expectErr: false,
		},
		{
			name:           "Invalid_Stress_Response",
			documentId:     "2",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Stress{},
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

			client := go_oura.NewClientWithUrlAndHttp("<<TOKEN>>", server.URL, server.Client())

			stress, err := client.GetStress(tc.documentId)
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

			if !reflect.DeepEqual(stress, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, stress)
			}
		})
	}
}

func TestGetStressDocuments(t *testing.T) {
	tt := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		mockResponse   string
		expectedOutput go_oura.Stresses
		expectErr      bool
	}{
		{
			name:         "Valid_Stresses_Response",
			startTime:    time.Now().Add(-1 * time.Hour),
			endTime:      time.Now().Add(-2 * time.Hour),
			mockResponse: `{"data":[{"id":"1304bda1-d5a0-4243-a321-0e48095149cc","day":"2024-01-09","stress_high":3600,"recovery_high":900,"day_summary":null},{"id":"1e08f663-2efa-41aa-901b-99b17cbeafcd","day":"2024-01-10","stress_high":900,"recovery_high":4500,"day_summary":"normal"}],"next_token":null}`,
			expectedOutput: go_oura.Stresses{
				Items: []go_oura.Stress{
					{
						ID: "1304bda1-d5a0-4243-a321-0e48095149cc",
						Day: func() go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2024-01-09")
							return go_oura.Date{Time: t}
						}(),
						StressHigh:   3600,
						RecoveryHigh: 900,
						DaySummary:   "",
					},
					{
						ID: "1e08f663-2efa-41aa-901b-99b17cbeafcd",
						Day: func() go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2024-01-10")
							return go_oura.Date{Time: t}
						}(),
						StressHigh:   900,
						RecoveryHigh: 4500,
						DaySummary:   "normal",
					},
				},
				NextToken: "",
			},
			expectErr: false,
		},
		{
			name:           "Invalid_Stresses_Response",
			startTime:      time.Now(),
			endTime:        time.Now(),
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Stresses{},
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

			client := go_oura.NewClientWithUrlAndHttp("<<TOKEN>>", server.URL, server.Client())

			stresses, err := client.GetStresses(tc.startTime, tc.endTime, nil)
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

			if !reflect.DeepEqual(stresses, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, stresses)
			}
		})
	}
}
