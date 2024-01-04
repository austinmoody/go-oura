package go_oura

type DailyReadiness struct {
	Id           string `json:"id"`
	Contributors struct {
		ActivityBalance     int `json:"activity_balance"`
		BodyTemperature     int `json:"body_temperature"`
		HrvBalance          int `json:"hrv_balance"`
		PreviousDayActivity int `json:"previous_day_activity"`
		PreviousNight       int `json:"previous_night"`
		RecoveryIndex       int `json:"recovery_index"`
		RestingHeartRate    int `json:"resting_heart_rate"`
		SleepBalance        int `json:"sleep_balance"`
	} `json:"contributors"`
	Day                       string  `json:"day"`
	Score                     int     `json:"score"`
	TemperatureDeviation      float64 `json:"temperature_deviation"`
	TemperatureTrendDeviation float64 `json:"temperature_trend_deviation"`
	Timestamp                 string  `json:"timestamp"`
}
