package go_oura

import "time"

type HeartRates struct {
	Items     []HeartRate `json:"data"`
	NextToken string      `json:"next_token"`
}

type HeartRate struct {
	Bpm       int       `json:"bpm"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

type heartRatesBase HeartRates
type hearRateBase HeartRate
