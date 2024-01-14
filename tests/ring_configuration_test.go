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

func TestGetRingConfiguration(t *testing.T) {
	tt := []struct {
		name           string
		documentId     string
		mockResponse   string
		expectedOutput go_oura.RingConfiguration
		expectErr      bool
	}{
		{
			name:         "Valid_RingConfiguration_Response",
			documentId:   "1",
			mockResponse: `{"id":"123","color":"black","design":"heritage","firmware_version":"2.9.24","hardware_type":"gen3","set_up_at":null,"size":9}`,
			expectedOutput: go_oura.RingConfiguration{
				ID:              "123",
				Color:           "black",
				Design:          "heritage",
				FirmwareVersion: "2.9.24",
				HardwareType:    "gen3",
				SetUpAt:         nil,
				Size:            9,
			},
			expectErr: false,
		},
		{
			name:           "Invalid_RingConfiguration_Response",
			documentId:     "2",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.RingConfiguration{},
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

			ringConfiguration, err := client.GetRingConfiguration(tc.documentId)
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

			if !reflect.DeepEqual(ringConfiguration, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, ringConfiguration)
			}
		})
	}
}

func TestGetRingConfigurations(t *testing.T) {
	tt := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		mockResponse   string
		expectedOutput go_oura.RingConfigurations
		expectErr      bool
	}{
		{
			name:         "Valid_RingConfigurations_Response",
			startTime:    time.Now().Add(-1 * time.Hour),
			endTime:      time.Now().Add(-2 * time.Hour),
			mockResponse: `{"data":[{"id":"123","color":"black","design":"round","firmware_version":"1.0","hardware_type":"Oura","set_up_at":"2022-04-01T00:00:00+00:00","size":9}],"next_token":null}`,
			expectedOutput: go_oura.RingConfigurations{
				Items: []go_oura.RingConfiguration{
					{
						ID:              "123",
						Color:           "black",
						Design:          "round",
						FirmwareVersion: "1.0",
						HardwareType:    "Oura",
						SetUpAt: func() *time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2022-04-01T00:00:00+00:00")
							return &t
						}(),
						Size: 9,
					},
				},
				NextToken: "",
			},
			expectErr: false,
		},
		{
			name:           "Invalid_RingConfigurations_Response",
			startTime:      time.Now().Add(-3 * time.Hour),
			endTime:        time.Now().Add(-4 * time.Hour),
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.RingConfigurations{},
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

			ringConfigurations, err := client.GetRingConfigurations(tc.startTime, tc.endTime, nil)
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

			if !reflect.DeepEqual(ringConfigurations, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, ringConfigurations)
			}
		})
	}
}
