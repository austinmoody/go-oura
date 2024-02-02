package go_oura

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	ouraApiUrlv2 = "https://api.ouraring.com/v2"
)

const (
	ActivityUrl          = "/usercollection/daily_activity"
	ReadinessUrl         = "/usercollection/daily_readiness"
	DailySleepUrl        = "/usercollection/daily_sleep"
	HeartRateUrl         = "/usercollection/heartrate"
	Spo2Url              = "/usercollection/daily_spo2"
	StressUrl            = "/usercollection/daily_stress"
	TagUrl               = "/usercollection/enhanced_tag"
	PersonalInfoUrl      = "/usercollection/personal_info"
	RestModeUrl          = "/usercollection/rest_mode_period"
	RingConfigurationUrl = "/usercollection/ring_configuration"
	SessionUrl           = "/usercollection/session"
	SleepUrl             = "/usercollection/sleep"
	SleepTimeUrl         = "/usercollection/sleep_time"
	WorkoutUrl           = "/usercollection/workout"
)

// Literally here so I can mock in tests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	accessToken string
	baseUrl     string
	HTTPClient  HTTPClient
}

func GetConfig(accessToken string) ClientConfig {
	return ClientConfig{
		accessToken: accessToken,
		baseUrl:     ouraApiUrlv2,
		HTTPClient:  &http.Client{},
	}
}

func GetConfigWithUrl(accessToken string, baseUrl string) ClientConfig {
	return ClientConfig{
		accessToken: accessToken,
		baseUrl:     baseUrl,
		HTTPClient:  &http.Client{},
	}
}

func GetConfigWithUrlAndHttp(accessToken string, baseUrl string, client HTTPClient) ClientConfig {
	return ClientConfig{
		accessToken: accessToken,
		baseUrl:     baseUrl,
		HTTPClient:  client,
	}
}

func (c *ClientConfig) GetUrl() (*url.URL, error) {
	apiUrl, err := url.ParseRequestURI(c.baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base url with error: %w", err)
	}

	return apiUrl, nil
}

func (c *ClientConfig) AddAuthorizationHeader(request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+c.accessToken)
}
