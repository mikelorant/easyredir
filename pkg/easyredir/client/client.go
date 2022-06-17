package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/mikelorant/easyredir/pkg/easyredir/option"
	"github.com/mikelorant/easyredir/pkg/jsonutil"
)

func New(opts ...option.Option) *Client {
	o := &option.Options{}

	for _, opt := range opts {
		opt.Apply(o)
	}

	return &Client{
		HTTPClient: buildHTTPClient(o),
		Config:     buildConfig(o),
	}
}

func (cl *Client) SendRequest(path, method string, body io.Reader) (io.ReadCloser, error) {
	url := fmt.Sprintf("%v%v", cl.Config.BaseURL, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new request: %w", err)
	}

	req.SetBasicAuth(cl.Config.APIKey, cl.Config.APISecret)
	req.Header.Set("Content-Type", ResourceType)
	req.Header.Set("Accept", ResourceType)

	if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" {
		req.Header.Set("Idempotency-Key", uuid.NewString())
	}

	resp, err := cl.HTTPClient.Do(req)
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

func buildConfig(opts *option.Options) *Config {
	cfg := &Config{
		BaseURL: BaseURL,
	}

	if opts.BaseURL != "" {
		cfg.BaseURL = opts.BaseURL
	}

	if opts.APIKey != "" {
		cfg.APIKey = opts.APIKey
	}

	if opts.APISecret != "" {
		cfg.APISecret = opts.APISecret
	}

	return cfg
}

func buildHTTPClient(opts *option.Options) Doer {
	if opts.HTTPClient != nil {
		return opts.HTTPClient
	}

	return http.DefaultClient
}
