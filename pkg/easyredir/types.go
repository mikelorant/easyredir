package easyredir

import "fmt"

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

type Metadata struct {
	HasMore bool `json:"has_more,omitempty"`
}

type Links struct {
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

type Pagination struct {
	StartingAfter string
	EndingBefore  string
}

type Options struct {
	SourceFilter string
	TargetFilter string
	Limit        int
	Pagination   Pagination
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
