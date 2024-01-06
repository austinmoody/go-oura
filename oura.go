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

func NewClient(accessToken string) *Client {

	return &Client{
		config: DefaultConfig(accessToken),
	}
}

func (c *Client) NewRequest(apiUrlPart string, params url.Values) (*http.Request, *OuraError) {
	apiUrl, err := url.Parse(c.config.BaseUrl)
	if err != nil {
		return nil,
			&OuraError{
				Code:    -1,
				Message: fmt.Sprintf("failed to parse base url with error: %v", err),
			}
	}

	apiUrl.Path = path.Join(apiUrl.Path, apiUrlPart)

	if params != nil && len(params) > 0 {
		apiUrl.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return nil,
			&OuraError{
				Code:    -1,
				Message: fmt.Sprintf("failed to create a new HTTP GET request with error: %v", err),
			}
	}
	req.Header.Set("Authorization", "Bearer "+c.config.accessToken)

	return req, nil
}

func (c *Client) Getter(apiUrlPart string, queryParams url.Values) (*[]byte, *OuraError) {

	req, ouraError := c.NewRequest(apiUrlPart, queryParams)
	if ouraError != nil {
		return nil,
			ouraError
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil,
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to complete HTTP request with error: %v", err),
			}
	}

	// Check for non-200 HTTP Status Code
	// TODO cater the message to HTTP Status Code per Oura API Spec
	if resp.StatusCode != 200 {
		return nil,
			&OuraError{
				Code:    resp.StatusCode,
				Message: fmt.Sprintf("received non-200 HTTP status code: %d", resp.StatusCode),
			}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil,
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to read response body with error: %v", err),
			}
	}

	err = resp.Body.Close()
	if err != nil {
		return nil,
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to close response body with error: %v", err),
			}
	}

	return &data, nil

}
