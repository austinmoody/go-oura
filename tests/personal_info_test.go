package tests

import (
	"github.com/austinmoody/go_oura"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetPersonalInfo(t *testing.T) {
	tt := []struct {
		name           string
		mockResponse   string
		expectedOutput go_oura.PersonalInfo
		expectErr      bool
	}{
		{
			name:         "Valid_PersonalInfo_Response",
			mockResponse: `{"id":"040aca77-8843-4221-8434-ee7d97d57ec9","age":25,"weight":82.3,"height":1.9,"biological_sex":"male","email":"whatever@example.com"}`,
			expectedOutput: go_oura.PersonalInfo{
				ID:     "040aca77-8843-4221-8434-ee7d97d57ec9",
				Age:    25,
				Weight: 82.3,
				Height: 1.9,
				Sex:    "male",
				Email:  "whatever@example.com",
			},
			expectErr: false,
		}, {
			name:           "Invalid_PersonalInfo_Response",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.PersonalInfo{},
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

			activity, err := client.GetPersonalInfo()
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
