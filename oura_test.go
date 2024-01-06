package go_oura

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

/*
Trying to convince myself these tests are useful... :)
*/
func TestNewClient(t *testing.T) {
	tt := []struct {
		name        string
		accessToken string
		expected    string
	}{
		{
			name:        "EmptyAccessToken",
			accessToken: "",
			expected:    "",
		},
		{
			name:        "ValidAccessToken",
			accessToken: "valid_token",
			expected:    "valid_token",
		},
		{
			name:        "AccessTokenWithSpecialChars",
			accessToken: "token#with$pecial@charact^rs",
			expected:    "token#with$pecial@charact^rs",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient(tc.accessToken)
			if client.config.accessToken != tc.expected {
				t.Errorf("NewClient(%q).accessToken = %q; expected %q", tc.accessToken, client.config.accessToken, tc.expected)
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
			// Create a new Client.
			client := &Client{
				config: ClientConfig{
					BaseUrl:     tc.baseUrl,
					accessToken: "test-token",
					HTTPClient:  &http.Client{},
				},
			}

			_, err := client.NewRequest(tc.apiUrlPart, tc.params)

			if tc.expectErr && err == nil {
				t.Errorf("expected an error but got none")
			}

			if !tc.expectErr && err != nil {
				t.Errorf("did not expect an error but got: %v", err)
			}

			// If there is an error, it must be OuraError.
			if err != nil {
				var ouraErr *OuraError
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
			// Create new client with mock HTTP client
			client := &Client{
				config: ClientConfig{
					BaseUrl:     "https://example.com",
					HTTPClient:  mockHTTPClient,
					accessToken: "accessToken",
				},
			}

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
