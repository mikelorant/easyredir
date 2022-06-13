package easyredir

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type Easyredir struct {
	Client ClientAPI
	Config *Config
}

type ClientAPI interface {
	SendRequest(baseURL, path, method string, body io.Reader) (io.ReadCloser, error)
}

type Client struct {
	httpClient *http.Client
	config     *Config
}

type Config struct {
	BaseURL   string
	APIKey    string
	APISecret string
}

const (
	_BaseURL      = "https://api.easyredir.com/v1"
	_ResourceType = "application/json; charset=utf-8"
)

func New(cfg *Config) *Easyredir {
	cfg.BaseURL = _BaseURL

	return &Easyredir{
		Client: &Client{
			httpClient: &http.Client{},
			config:     cfg,
		},
		Config: cfg,
	}
}

func (c *Easyredir) Ping() string {
	return "pong"
}

func (cl *Client) SendRequest(baseURL, path, method string, body io.Reader) (io.ReadCloser, error) {
	url := fmt.Sprintf("%v%v", baseURL, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new request: %w", err)
	}

	req.SetBasicAuth(cl.config.APIKey, cl.config.APISecret)
	req.Header.Set("Content-Type", _ResourceType)
	req.Header.Set("Accept", _ResourceType)

	if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" {
		req.Header.Set("Idempotency-Key", uuid.NewString())
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to do request: %w", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, &RateLimitError{
			Limit:     resp.Header.Get("X-Ratelimit-Limit"),
			Remaining: resp.Header.Get("X-Ratelimit-Remaining"),
			Reset:     resp.Header.Get("X-Ratelimit-Reset"),
		}
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		apiErr := APIErrors{}
		if err := decodeJSON(resp.Body, &apiErr); err == nil {
			return nil, apiErr
		}
		return nil, fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func decodeJSON(r io.ReadCloser, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return fmt.Errorf("unable to json decode: %w", err)
	}
	r.Close()

	return nil
}
