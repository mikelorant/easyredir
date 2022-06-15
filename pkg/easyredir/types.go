package easyredir

import (
	"fmt"
	"net/http"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir/option"
)

type Option interface {
	Apply(*option.Options)
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
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
