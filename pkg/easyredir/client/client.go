package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/config"
	"github.com/mikelorant/easyredir-cli/pkg/jsonutil"
)

type Client struct {
	httpClient *http.Client
	username   string
	password   string
}

const (
	_ResourceType = "application/json; charset=utf-8"
)

func New(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{},
		username: cfg.APIKey(),
		password: cfg.APISecret(),
	}
}

func (cl *Client) SendRequest(baseURL, path, method string, body io.Reader) (io.ReadCloser, error) {
	url := fmt.Sprintf("%v%v", baseURL, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new request: %w", err)
	}

	req.SetBasicAuth(cl.username, cl.password)
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
		if err := jsonutil.DecodeJSON(resp.Body, &apiErr); err == nil {
			return nil, apiErr
		}
		return nil, fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
