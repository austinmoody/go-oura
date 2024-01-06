package go_oura

import "net/http"

const (
	ouraApiUrlv2 = "https://api.ouraring.com/v2"
)

// Literally here so I can mock in tests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	accessToken string
	BaseUrl     string
	HTTPClient  HTTPClient
}

func DefaultConfig(accessToken string) ClientConfig {
	return ClientConfig{
		accessToken: accessToken,
		BaseUrl:     ouraApiUrlv2,
		HTTPClient:  &http.Client{},
	}
}
