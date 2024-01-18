package go_oura

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	Config ClientConfig
}

/*
NewClient

Function Name:

	NewClient(accessToken string) *Client

Description:

	The 'NewClient' function is used to create a new go_oura client with given access token.

	The Client will use the Base URL specified in config.go.

Parameters:
 1. accessToken: A string representing your Oura Ring personal access token required to authenticate the client.

Returns:

	A pointer to a new 'Client' structure instance.

Example usage:

	client := NewClient("your_access_token_here")
*/
func NewClient(accessToken string) *Client {
	return &Client{
		Config: GetConfig(accessToken),
	}
}

/*
NewClientWithUrl

Function Name:

	NewClientWithUrl(accessToken string, baseUrl string) *Client

Description:

	The 'NewClientWithUrl' function is used to create a new go_oura client with given access token and base url.

Parameters:

 1. accessToken: A string representing your Oura Ring personal access token required to authenticate the client.
 2. baseUrl: The base API url for the Oura Ring API.  The default, for example, is https://api.ouraring.com/v2

Returns:

	A pointer to a new 'Client' structure instance.

Example usage:

	client := NewClientWithUrl("your_access_token_here", "https://api.ouraring.com/v2")
*/
func NewClientWithUrl(accessToken string, baseUrl string) *Client {
	return &Client{
		Config: GetConfigWithUrl(accessToken, baseUrl),
	}
}

func NewClientWithUrlAndHttp(accessToken string, baseUrl string, client HTTPClient) *Client {
	return &Client{
		Config: GetConfigWithUrlAndHttp(accessToken, baseUrl, client),
	}
}

func (c *Client) NewRequest(apiUrlPart string, params url.Values) (*http.Request, *OuraError) {

	apiUrl, ouraError := c.Config.GetUrl()
	if ouraError != nil {
		return nil,
			ouraError
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
	c.Config.AddAuthorizationHeader(req)

	return req, nil
}

func (c *Client) Getter(apiUrlPart string, queryParams url.Values) (*[]byte, *OuraError) {

	req, ouraError := c.NewRequest(apiUrlPart, queryParams)
	if ouraError != nil {
		return nil,
			ouraError
	}

	resp, err := c.Config.HTTPClient.Do(req)
	if err != nil {
		return nil,
			&OuraError{
				Code:    0,
				Message: fmt.Sprintf("failed to complete HTTP request with error: %v", err),
			}
	}

	// Check for non-200 HTTP Status Code
	if resp.StatusCode != 200 {
		return nil,
			&OuraError{
				Code:    resp.StatusCode,
				Message: resp.Status,
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
