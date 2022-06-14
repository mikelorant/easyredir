package option

import (
	"net/http"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir/pagination"
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
	Pagination   pagination.Pagination
}

func WithLimit(limit int) func(*Options) {
	return func(o *Options) {
		o.Limit = limit
	}
}

func WithSourceFilter(url string) func(*Options) {
	return func(o *Options) {
		o.SourceFilter = url
	}
}

func WithTargetFilter(url string) func(*Options) {
	return func(o *Options) {
		o.TargetFilter = url
	}
}

func WithAPIKey(apiKey string) func(*Options) {
	return func(o *Options) {
		o.APIKey = apiKey
	}
}

func WithAPISecret(apiSecret string) func(*Options) {
	return func(o *Options) {
		o.APISecret = apiSecret
	}
}

func WithBaseURL(baseURL string) func(*Options) {
	return func(o *Options) {
		o.BaseURL = baseURL
	}
}

func WithHTTPClient(client Doer) func(*Options) {
	return func(o *Options) {
		o.HTTPClient = client
	}
}
