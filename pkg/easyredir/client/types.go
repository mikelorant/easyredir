package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mikelorant/easyredir/pkg/structutil"
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
	var sb strings.Builder

	fmt.Fprint(&sb, err.Type)
	if err.Message != "" {
		fmt.Fprintf(&sb, ": %v", err.Message)
	}
	if len(err.Errors) > 0 {
		ae, _ := structutil.Sprint(err.Errors)
		fmt.Fprintf(&sb, "\nerrors:\n%v", ae)
	}

	return sb.String()
}

func (err RateLimitError) Error() string {
	return fmt.Sprintf("rate limited with limit: %v, remaining: %v, reset: %v", err.Limit, err.Remaining, err.Reset)
}
