package easyredir

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Easyredir struct {
	client ClientAPI
	config *Config
}

type ClientAPI interface {
	sendRequest(baseURL, path, method string, body io.Reader) (io.ReadCloser, error)
}

type Client struct {
	httpClient *http.Client
	config     *Config
}

type Config struct {
	baseURL   string
	APIKey    string
	APISecret string
}

type APIErrors struct {
	Type    string     `json:"type"`
	Message string     `json:"message"`
	Errors  []APIError `json:"errors"`
}

type APIError struct {
	Resource string `json:"resource"`
	Param    string `json:"param"`
	Code     string `json:"code"`
	Message  string `json:"message"`
}

const (
	_BaseURL      = "https://api.easyredir.com/v1"
	_ResourceType = "application/json; charset=utf-8"
)

func New(cfg *Config) *Easyredir {
	cfg.baseURL = _BaseURL

	return &Easyredir{
		client: &Client{
			httpClient: &http.Client{},
			config:     cfg,
		},
		config: cfg,
	}
}

func (c *Easyredir) Ping() string {
	return "pong"
}

func (err APIErrors) Error() string {
	var sb strings.Builder
	fmt.Fprint(&sb, err.Type)
	if err.Message != "" {
		fmt.Fprintf(&sb, " :%v", err.Message)
	}
	return sb.String()
}

func decodeJSON(r io.ReadCloser, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return fmt.Errorf("unable to json decode: %w", err)
	}
	r.Close()

	return nil
}

func (cl *Client) sendRequest(baseURL, path, method string, body io.Reader) (io.ReadCloser, error) {
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
		return nil, fmt.Errorf("rate limited: remaining: %v, limit: %v, reset: %v",
			resp.Header.Get("X-Ratelimit-Limit"),
			resp.Header.Get("X-Ratelimit-Remaining"),
			resp.Header.Get("X-Ratelimit-Reset"),
		)
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
