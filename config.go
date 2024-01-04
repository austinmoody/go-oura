package go_oura

import "net/http"

const (
	ouraApiUrlv2 = "https://api.ouraring.com/v2"
)

type ClientConfig struct {
	accessToken string
	BaseUrl     string
	HTTPClient  *http.Client
}

func DefaultConfig(accessToken string) ClientConfig {
	return ClientConfig{
		accessToken: accessToken,
		BaseUrl:     ouraApiUrlv2,
		HTTPClient:  &http.Client{},
	}
}
