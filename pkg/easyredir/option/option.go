package option

import (
	"net/http"
)

type Option interface {
	Apply(*Options)
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Options struct {
	BaseURL      string
	APIKey       string
	APISecret    string
	HTTPClient   Doer
	SourceFilter string
	TargetFilter string
	Limit        int
	Include      string
	Pagination   Pagination
}
