package easyredir

import (
	"fmt"
	"net/http"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	HTTPClient Doer
	Config     *Config
}

type Config struct {
	BaseURL   string
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

type RateLimitError struct {
	Limit     string
	Remaining string
	Reset     string
}

const (
	BaseURL = "https://api.easyredir.com/v1"
)

const (
	ResourceType = "application/json; charset=utf-8"
)

func (err APIErrors) Error() string {
	str := err.Type
	if err.Message != "" {
		str = fmt.Sprintf("%v: %v", str, err.Message)
	}
	return str
}

func (err RateLimitError) Error() string {
	return fmt.Sprintf("rate limited with limit: %v, remaining: %v, reset: %v", err.Limit, err.Remaining, err.Reset)
}
