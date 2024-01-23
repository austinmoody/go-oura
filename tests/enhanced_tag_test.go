package tests

import (
	"errors"
	"github.com/austinmoody/go_oura"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetTag(t *testing.T) {
	tt := []struct {
		name           string
		documentId     string
		mockResponse   string
		expectedOutput go_oura.EnhancedTag
		expectErr      bool
	}{
		{
			name:         "Valid_Tag_Response",
			documentId:   "1",
			mockResponse: `{"id":"f70f8344-8ec0-47b3-a4d8-b610432cb722","tag_type_code":"tag_generic_anxiety","start_time":"2023-12-21T21:51:31-05:00","end_time":"2023-12-21T23:51:31-05:00","start_day":"2023-12-21","end_day":"2023-12-21","comment":"This is a comment"}`,
			expectedOutput: go_oura.EnhancedTag{
				ID:          "f70f8344-8ec0-47b3-a4d8-b610432cb722",
				TagTypeCode: "tag_generic_anxiety",
				StartTime: func() *time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2023-12-21T21:51:31-05:00")
					return &t
				}(),
				EndTime: func() *time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2023-12-21T23:51:31-05:00")
					return &t
				}(),
				StartDay: func() *go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2023-12-21")
					return &go_oura.Date{Time: t}
				}(),
				EndDay: func() *go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2023-12-21")
					return &go_oura.Date{Time: t}
				}(),
				Comment: "This is a comment",
			},
			expectErr: false,
		},
		{
			name:           "Invalid_Tag_Response",
			documentId:     "2",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.EnhancedTag{},
			expectErr:      true,
		},
		{
			name:         "MissingFields_Tag_Response",
			documentId:   "3",
			mockResponse: `{"id":"cfe144cf-c954-4609-b3d2-aa3291cd3dd3","tag_type_code":"tag_generic_nocaffeine","start_time":"2023-12-04T18:01:00-05:00","end_time":null,"start_day":"2023-12-04","end_day":null,"comment":null}`,
			expectedOutput: go_oura.EnhancedTag{
				ID:          "cfe144cf-c954-4609-b3d2-aa3291cd3dd3",
				TagTypeCode: "tag_generic_nocaffeine",
				StartTime: func() *time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2023-12-04T18:01:00-05:00")
					return &t
				}(),
				EndTime: nil,
				StartDay: func() *go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2023-12-04")
					return &go_oura.Date{Time: t}
				}(),
				EndDay:  nil,
				Comment: "",
			},
			expectErr: false,
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

			tag, err := client.GetEnhancedTag(tc.documentId)
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

			if !reflect.DeepEqual(tag, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, tag)
			}
		})
	}
}

func TestGetTagDocuments(t *testing.T) {
	tt := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		mockResponse   string
		expectedOutput go_oura.EnhancedTags
		expectErr      bool
	}{
		{
			name:         "Valid_Tags_Response",
			startTime:    time.Now().Add(-1 * time.Hour),
			endTime:      time.Now().Add(-2 * time.Hour),
			mockResponse: `{"data":[{"id":"f70f8344-8ec0-47b3-a4d8-b610432cb722","tag_type_code":"tag_generic_anxiety","start_time":"2023-12-21T21:51:31-05:00","end_time":"2023-12-21T23:51:31-05:00","start_day":"2023-12-21","end_day":"2023-12-21","comment":"This is a comment"}],"next_token":null}`,
			expectedOutput: go_oura.EnhancedTags{
				Items: []go_oura.EnhancedTag{
					{
						ID:          "f70f8344-8ec0-47b3-a4d8-b610432cb722",
						TagTypeCode: "tag_generic_anxiety",
						StartTime: func() *time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2023-12-21T21:51:31-05:00")
							return &t
						}(),
						EndTime: func() *time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2023-12-21T23:51:31-05:00")
							return &t
						}(),
						StartDay: func() *go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2023-12-21")
							return &go_oura.Date{Time: t}
						}(),
						EndDay: func() *go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2023-12-21")
							return &go_oura.Date{Time: t}
						}(),
						Comment: "This is a comment",
					},
				},
				NextToken: "",
			},
			expectErr: false,
		}, {
			name:           "Invalid_TagDocuments_Response",
			startTime:      time.Now().Add(-3 * time.Hour),
			endTime:        time.Now().Add(-4 * time.Hour),
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.EnhancedTags{},
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

			tags, err := client.GetEnhancedTags(tc.startTime, tc.endTime, nil)
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

			if !reflect.DeepEqual(tags, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, tags)
			}
		})
	}
}
