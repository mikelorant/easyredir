package option

import (
	"net/http"
)

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
	Pagination   Pagination
}
