package go_oura

type Client struct {
	config ClientConfig
}

func NewClient(accessToken string) *Client {

	return &Client{
		config: DefaultConfig(accessToken),
	}
}
