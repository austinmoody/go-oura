package tests

import (
	"errors"
	"fmt"
	"github.com/austinmoody/go-oura"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

/*
Trying to convince myself these tests are useful... :)
*/
func TestNewClientWithUrl(t *testing.T) {
	tt := []struct {
		name      string
		baseUrl   string
		expected  string
		expectErr bool
	}{
		{
			name:      "ValidBaseUrl",
			baseUrl:   "http://www.example.com/",
			expected:  "http://www.example.com/",
			expectErr: false,
		},
		{
			name:      "InvalidUrl",
			baseUrl:   "This Is Not A Url",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "BlankUrl",
			baseUrl:   "",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "ActualOuraUrl",
			baseUrl:   "https://api.ouraring.com/v2/usercollection/daily_activity/45173cbe-ef26-430f-adc4-c4a1424b45ab",
			expected:  "https://api.ouraring.com/v2/usercollection/daily_activity/45173cbe-ef26-430f-adc4-c4a1424b45ab",
			expectErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := go_oura.NewClientWithUrl("", tc.baseUrl)

			apiUrl, err := client.Config.GetUrl()

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

			if apiUrl.String() != tc.expected {
				t.Errorf("NewClientWithUrl created with BaseUrl %q, expected %q got %q", tc.baseUrl, tc.baseUrl, apiUrl.String())
			}
		})
	}
}

func TestClient_NewRequest(t *testing.T) {

	testCases := []struct {
		name       string
		baseUrl    string
		apiUrlPart string
		params     url.Values
		expectErr  bool
	}{
		{
			name:       "valid_parameters",
			baseUrl:    "http://localhost:8080",
			apiUrlPart: "/api/user",
			params:     url.Values{"key": []string{"value"}},
			expectErr:  false,
		},
		{
			name:       "empty_params",
			baseUrl:    "http://localhost:8080",
			apiUrlPart: "/api/user",
			params:     nil,
			expectErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			clientConfig := go_oura.GetConfigWithUrl("test-token", tc.baseUrl)

			// Create a new Client.
			client := &go_oura.Client{Config: clientConfig}

			_, err := client.NewRequest(tc.apiUrlPart, tc.params)

			if tc.expectErr && err == nil {
				t.Errorf("expected an error but got none")
			}

			if !tc.expectErr && err != nil {
				t.Errorf("did not expect an error but got: %v", err)
			}

			// If there is an error, it must be OuraError.
			if err != nil {
				var ouraErr *go_oura.OuraError
				if !errors.As(err, &ouraErr) {
					t.Errorf("expected an OuraError but got a different error: %v", err)
				}
			}
		})
	}
}

func TestClient_Getter(t *testing.T) {
	mockHTTPClient := NewMockHTTPClient()

	tests := []struct {
		name         string
		apiUrl       string
		queryParams  url.Values
		mockResponse *http.Response
		mockErr      error
		wantErr      bool
	}{
		{
			name:        "successful request",
			apiUrl:      "/usercollection/daily_readiness",
			queryParams: url.Values{"start_date": {"2006-01-02"}},
			mockResponse: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success"}`)),
			},
			wantErr: false,
		},
		{
			name:        "client send error",
			apiUrl:      "/usercollection/daily_readiness",
			queryParams: url.Values{"start_date": {"2006-01-02"}},
			mockErr:     fmt.Errorf("mock error"),
			wantErr:     true,
		},
		{
			name:        "non-successful status code",
			apiUrl:      "/usercollection/daily_readiness",
			queryParams: url.Values{"start_date": {"2006-01-02"}},
			mockResponse: &http.Response{
				StatusCode: 400,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clientConfig := go_oura.GetConfigWithUrlAndHttp("accessToken", "", mockHTTPClient)

			// Create new client with mock HTTP client
			client := &go_oura.Client{Config: clientConfig}

			mockHTTPClient.NextResponse = tc.mockResponse
			mockHTTPClient.NextErr = tc.mockErr

			_, err := client.Getter(tc.apiUrl, tc.queryParams)

			if (err != nil) != tc.wantErr {
				t.Errorf("Client.Get() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

type MockHTTPClient struct {
	NextResponse *http.Response
	NextErr      error
}

func (c *MockHTTPClient) Do(*http.Request) (*http.Response, error) {
	return c.NextResponse, c.NextErr
}

func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{}
}
