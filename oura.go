package go_oura

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	config ClientConfig
}

type OuraApiResponse struct {
	Code int
	Body []byte
}

func NewClient(accessToken string) *Client {

	return &Client{
		config: DefaultConfig(accessToken),
	}
}

func (c *Client) NewRequest(apiUrlPart string, params url.Values) (*http.Request, error) {
	apiUrl, err := url.Parse(c.config.BaseUrl)
	if err != nil {
		return nil,
			fmt.Errorf("failed to parse base url with error: %w", err)
	}

	apiUrl.Path = path.Join(apiUrl.Path, apiUrlPart)

	if params != nil && len(params) > 0 {
		apiUrl.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return nil,
			fmt.Errorf("failed to create a new HTTP GET request with error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.config.accessToken)

	return req, nil
}

func (c *Client) Getter(apiUrlPart string, queryParams url.Values) (OuraApiResponse, error) {

	req, err := c.NewRequest(apiUrlPart, queryParams)
	if err != nil {
		return OuraApiResponse{},
			fmt.Errorf("failed to create new Client Getter with error: %w", err)
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return OuraApiResponse{},
			fmt.Errorf("failed to complete HTTP request with error: %w", err)
	}

	apiResponse := OuraApiResponse{resp.StatusCode, nil}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiResponse,
			fmt.Errorf("failed to read response body with error: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return apiResponse,
			fmt.Errorf("failed to close response body with error: %w", err)
	}

	apiResponse.Body = data

	return apiResponse, nil

}
