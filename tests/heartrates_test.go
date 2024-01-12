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

func TestGetHeartRates(t *testing.T) {
	tt := []struct {
		name           string
		startDateTime  time.Time
		endDateTime    time.Time
		mockResponse   string
		expectedOutput go_oura.HeartRates
		expectErr      bool
	}{
		{
			name:          "Valid_HeartRates_Response",
			startDateTime: time.Now().Add(-1 * time.Hour),
			endDateTime:   time.Now().Add(-2 * time.Hour),
			mockResponse:  `{"data":[{"bpm":74,"source":"awake","timestamp":"2024-01-10T01:45:45+00:00"},{"bpm":95,"source":"awake","timestamp":"2024-01-10T01:46:00+00:00"}],"next_token":"1:eyJ1Ijp7IjkwMTJhYjcwLTY4YWMtNGRkYi04ZGI2LTk5MTA3Mjk2NTMwYSI6eyJuZXdfc3RhcnRfdGltZSI6IjIwMjQtMDEtMDZUMjI6NDQ6NTEuODAxWiJ9fSwibiI6W119"}`,
			expectedOutput: go_oura.HeartRates{
				Items: []go_oura.HeartRate{
					{
						Bpm:    74,
						Source: "awake",
						Timestamp: func() time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2024-01-10T01:45:45+00:00")
							return t
						}(),
					}, {
						Bpm:    95,
						Source: "awake",
						Timestamp: func() time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2024-01-10T01:46:00+00:00")
							return t
						}(),
					},
				},
				NextToken: "1:eyJ1Ijp7IjkwMTJhYjcwLTY4YWMtNGRkYi04ZGI2LTk5MTA3Mjk2NTMwYSI6eyJuZXdfc3RhcnRfdGltZSI6IjIwMjQtMDEtMDZUMjI6NDQ6NTEuODAxWiJ9fSwibiI6W119",
			},
			expectErr: false,
		},
		{
			name:           "Invalid_HeartRates_Response",
			startDateTime:  time.Now().Add(-3 * time.Hour),
			endDateTime:    time.Now().Add(-4 * time.Hour),
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.HeartRates{},
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

			activity, err := client.GetHeartRates(tc.startDateTime, tc.endDateTime, nil)
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
