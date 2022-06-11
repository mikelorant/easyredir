package easyredir

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Easyredir struct {
	client ClientAPI
	config *Config
}

type ClientAPI interface {
	sendRequest(baseURL, path, method string, body io.Reader) (io.Reader, error)
}

type Client struct {
	httpClient *http.Client
}

type Config struct {
	BaseURL	string
	Key    	string
	Secret 	string
}

func New(cfg *Config) *Easyredir {
	return &Easyredir{
		client: &Client{
			httpClient: &http.Client{},
		},
		config: cfg,
	}
}

func (c *Easyredir) Ping() string {
	return "pong"
}

func decodeJSON(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return fmt.Errorf("unable to json decode: %w", err)
	}

	return nil
}

func (cl *Client) sendRequest(baseURL, path, method string, body io.Reader) (io.Reader, error) {
	url := fmt.Sprintf("%v%v", baseURL, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new request: %w", err)
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to do request: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
