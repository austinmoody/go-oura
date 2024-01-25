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

func TestGetSleep(t *testing.T) {
	tt := []struct {
		name           string
		sleepId        string
		mockResponse   string
		expectedOutput go_oura.Sleep
		expectErr      bool
	}{
		{
			name:         "Valid_Sleep_Document",
			sleepId:      "1",
			mockResponse: `{"id":"02e54fa5-f514-46f6-9ac6-b651d8f84b75","average_breath":13.25,"average_heart_rate":71.5,"average_hrv":17,"awake_time":6650,"bedtime_end":"2024-01-21T09:08:18-05:00","bedtime_start":"2024-01-21T01:28:28-05:00","day":"2024-01-21","deep_sleep_duration":2850,"efficiency":76,"heart_rate":{"interval":300,"items":[64,null],"timestamp":"2024-01-21T01:28:28.000-05:00"},"hrv":{"interval":300,"items":[30,13,null,42,null,15],"timestamp":"2024-01-21T01:28:28.000-05:00"},"latency":180,"light_sleep_duration":12690,"low_battery_alert":false,"lowest_heart_rate":64,"movement_30_sec":"11111111111111211111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111121111111111111111133223444222222244444444332333221111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111122111111111111111111111111111111111111111111111111111111111111111111321111111111111111111111111111111111112111111111113432222233233231343111111111111111221111111111111111111211111111111111121111111111111111111111111111111111111111111111111111111111111111111111111111211111111111221111111111111111111111111111111111111111111111333323332111111121111111334323332333323343231111111111111111111111111111111211111111111111111111111111111111111111111111111111111221111111111111111111111221111111111111111111113211332111111133","period":1,"readiness":{"contributors":{"activity_balance":74,"body_temperature":100,"hrv_balance":55,"previous_day_activity":91,"previous_night":41,"recovery_index":100,"resting_heart_rate":1,"sleep_balance":79},"score":61,"temperature_deviation":-0.09,"temperature_trend_deviation":-0.02},"readiness_score_delta":0,"rem_sleep_duration":5400,"restless_periods":32,"sleep_phase_5_min":"42222221222221122222444412211112212233332232222223344444424222222334343334444444422223222334","sleep_score_delta":0,"sleep_algorithm_version":"v2","time_in_bed":27590,"total_sleep_duration":20940,"type":"long_sleep"}`,
			expectedOutput: go_oura.Sleep{
				ID:               "02e54fa5-f514-46f6-9ac6-b651d8f84b75",
				AverageBreath:    13.25,
				AverageHeartRate: 71.5,
				AverageHrv:       17,
				AwakeTime:        6650,
				BedtimeEnd: func() time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2024-01-21T09:08:18-05:00")
					return t
				}(),
				BedtimeStart: func() time.Time {
					layout := "2006-01-02T15:04:05Z07:00"
					t, _ := time.Parse(layout, "2024-01-21T01:28:28-05:00")
					return t
				}(),
				Day: func() go_oura.Date {
					layout := "2006-01-02"
					t, _ := time.Parse(layout, "2024-01-21")
					return go_oura.Date{Time: t}
				}(),
				DeepSleepDuration:     2850,
				Efficiency:            76,
				Latency:               180,
				LightSleepDuration:    12690,
				LowBatteryAlert:       false,
				LowestHeartRate:       64,
				Movement30Sec:         "11111111111111211111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111121111111111111111133223444222222244444444332333221111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111122111111111111111111111111111111111111111111111111111111111111111111321111111111111111111111111111111111112111111111113432222233233231343111111111111111221111111111111111111211111111111111121111111111111111111111111111111111111111111111111111111111111111111111111111211111111111221111111111111111111111111111111111111111111111333323332111111121111111334323332333323343231111111111111111111111111111111211111111111111111111111111111111111111111111111111111221111111111111111111111221111111111111111111113211332111111133",
				Period:                1,
				ReadinessScoreDelta:   0,
				RemSleepDuration:      5400,
				RestlessPeriods:       32,
				SleepPhase5Min:        "42222221222221122222444412211112212233332232222223344444424222222334343334444444422223222334",
				SleepScoreDelta:       0,
				SleepAlgorithmVersion: "v2",
				TimeInBed:             27590,
				TotalSleepDuration:    20940,
				Type:                  "long_sleep",
				HeartRate: go_oura.IntervalItems{
					Interval: 300,
					Timestamp: func() time.Time {
						layout := "2006-01-02T15:04:05Z07:00"
						t, _ := time.Parse(layout, "2024-01-21T01:28:28.000-05:00")
						return t
					}(),
					Items: []float64{
						64,
						0,
					},
				},
				Hrv: go_oura.IntervalItems{
					Interval: 300,
					Timestamp: func() time.Time {
						layout := "2006-01-02T15:04:05Z07:00"
						t, _ := time.Parse(layout, "2024-01-21T01:28:28.000-05:00")
						return t
					}(),
					Items: []float64{
						30,
						13,
						0,
						42,
						0,
						15,
					},
				},
				Readiness: go_oura.SleepReadiness{
					Contributors: go_oura.Contributors{
						ActivityBalance:     74,
						BodyTemperature:     100,
						HrvBalance:          55,
						PreviousDayActivity: 91,
						PreviousNight:       41,
						RecoveryIndex:       100,
						RestingHeartRate:    1,
						SleepBalance:        79,
					},
					Score:                     61,
					TemperatureDeviation:      -0.09,
					TemperatureTrendDeviation: -0.02,
				},
			},
			expectErr: false,
		}, {
			name:           "Invalid_Sleep_Document",
			sleepId:        "2",
			mockResponse:   `{"message": "invalid"}`,
			expectedOutput: go_oura.Sleep{},
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

			activity, err := client.GetSleep(tc.sleepId)
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

func TestGetSleeps(t *testing.T) {
	tt := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		mockResponse   string
		expectedOutput go_oura.Sleeps
		expectErr      bool
	}{
		{
			name:         "Valid_Sleeps_Document",
			startTime:    time.Now().Add(-1 * time.Hour),
			endTime:      time.Now().Add(-2 * time.Hour),
			mockResponse: `{"data":[{"id":"02e54fa5-f514-46f6-9ac6-b651d8f84b75","average_breath":13.25,"average_heart_rate":71.5,"average_hrv":17,"awake_time":6650,"bedtime_end":"2024-01-21T09:08:18-05:00","bedtime_start":"2024-01-21T01:28:28-05:00","day":"2024-01-21","deep_sleep_duration":2850,"efficiency":76,"heart_rate":{"interval":300,"items":[64,null],"timestamp":"2024-01-21T01:28:28.000-05:00"},"hrv":{"interval":300,"items":[30,13,null,42,null,15],"timestamp":"2024-01-21T01:28:28.000-05:00"},"latency":180,"light_sleep_duration":12690,"low_battery_alert":false,"lowest_heart_rate":64,"movement_30_sec":"11111111111111211111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111121111111111111111133223444222222244444444332333221111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111122111111111111111111111111111111111111111111111111111111111111111111321111111111111111111111111111111111112111111111113432222233233231343111111111111111221111111111111111111211111111111111121111111111111111111111111111111111111111111111111111111111111111111111111111211111111111221111111111111111111111111111111111111111111111333323332111111121111111334323332333323343231111111111111111111111111111111211111111111111111111111111111111111111111111111111111221111111111111111111111221111111111111111111113211332111111133","period":1,"readiness":{"contributors":{"activity_balance":74,"body_temperature":100,"hrv_balance":55,"previous_day_activity":91,"previous_night":41,"recovery_index":100,"resting_heart_rate":1,"sleep_balance":79},"score":61,"temperature_deviation":-0.09,"temperature_trend_deviation":-0.02},"readiness_score_delta":0,"rem_sleep_duration":5400,"restless_periods":32,"sleep_phase_5_min":"42222221222221122222444412211112212233332232222223344444424222222334343334444444422223222334","sleep_score_delta":0,"sleep_algorithm_version":"v2","time_in_bed":27590,"total_sleep_duration":20940,"type":"long_sleep"}],"next_token":"your-token"}`,
			expectedOutput: go_oura.Sleeps{
				NextToken: "your-token",
				Items: []go_oura.Sleep{
					{ID: "02e54fa5-f514-46f6-9ac6-b651d8f84b75",
						AverageBreath:    13.25,
						AverageHeartRate: 71.5,
						AverageHrv:       17,
						AwakeTime:        6650,
						BedtimeEnd: func() time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2024-01-21T09:08:18-05:00")
							return t
						}(),
						BedtimeStart: func() time.Time {
							layout := "2006-01-02T15:04:05Z07:00"
							t, _ := time.Parse(layout, "2024-01-21T01:28:28-05:00")
							return t
						}(),
						Day: func() go_oura.Date {
							layout := "2006-01-02"
							t, _ := time.Parse(layout, "2024-01-21")
							return go_oura.Date{Time: t}
						}(),
						DeepSleepDuration:     2850,
						Efficiency:            76,
						Latency:               180,
						LightSleepDuration:    12690,
						LowBatteryAlert:       false,
						LowestHeartRate:       64,
						Movement30Sec:         "11111111111111211111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111121111111111111111133223444222222244444444332333221111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111122111111111111111111111111111111111111111111111111111111111111111111321111111111111111111111111111111111112111111111113432222233233231343111111111111111221111111111111111111211111111111111121111111111111111111111111111111111111111111111111111111111111111111111111111211111111111221111111111111111111111111111111111111111111111333323332111111121111111334323332333323343231111111111111111111111111111111211111111111111111111111111111111111111111111111111111221111111111111111111111221111111111111111111113211332111111133",
						Period:                1,
						ReadinessScoreDelta:   0,
						RemSleepDuration:      5400,
						RestlessPeriods:       32,
						SleepPhase5Min:        "42222221222221122222444412211112212233332232222223344444424222222334343334444444422223222334",
						SleepScoreDelta:       0,
						SleepAlgorithmVersion: "v2",
						TimeInBed:             27590,
						TotalSleepDuration:    20940,
						Type:                  "long_sleep",
						HeartRate: go_oura.IntervalItems{
							Interval: 300,
							Timestamp: func() time.Time {
								layout := "2006-01-02T15:04:05Z07:00"
								t, _ := time.Parse(layout, "2024-01-21T01:28:28.000-05:00")
								return t
							}(),
							Items: []float64{
								64,
								0,
							},
						},
						Hrv: go_oura.IntervalItems{
							Interval: 300,
							Timestamp: func() time.Time {
								layout := "2006-01-02T15:04:05Z07:00"
								t, _ := time.Parse(layout, "2024-01-21T01:28:28.000-05:00")
								return t
							}(),
							Items: []float64{
								30,
								13,
								0,
								42,
								0,
								15,
							},
						},
						Readiness: go_oura.SleepReadiness{
							Contributors: go_oura.Contributors{
								ActivityBalance:     74,
								BodyTemperature:     100,
								HrvBalance:          55,
								PreviousDayActivity: 91,
								PreviousNight:       41,
								RecoveryIndex:       100,
								RestingHeartRate:    1,
								SleepBalance:        79,
							},
							Score:                     61,
							TemperatureDeviation:      -0.09,
							TemperatureTrendDeviation: -0.02},
					},
				},
			},
			expectErr: false,
		}, {
			name:           "Invalid_Sleeps_Document",
			startTime:      time.Now().Add(-3 * time.Hour),
			endTime:        time.Now().Add(-4 * time.Hour),
			mockResponse:   `{"message": "invalid"}`,
			expectErr:      true,
			expectedOutput: go_oura.Sleeps{},
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

			activity, err := client.GetSleeps(tc.startTime, tc.endTime, nil)
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
