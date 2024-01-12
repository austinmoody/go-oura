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
	ActivityUrl  = "/usercollection/daily_activity"
	ReadinessUrl = "/usercollection/daily_readiness"
	SleepUrl     = "/usercollection/daily_sleep"
	HeartRateUrl = "/usercollection/heartrate"
	Spo2Url      = "/usercollection/daily_spo2"
	StressUrl    = "/usercollection/daily_stress"
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

func (c *ClientConfig) GetUrl() (*url.URL, *OuraError) {
	apiUrl, err := url.ParseRequestURI(c.baseUrl)
	if err != nil {
		return nil,
			&OuraError{
				Code:    -1,
				Message: fmt.Sprintf("failed to parse base url with error: %v", err),
			}
	}

	return apiUrl, nil
}

func (c *ClientConfig) AddAuthorizationHeader(request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+c.accessToken)
}
