package option

import "github.com/mikelorant/easyredir-cli/pkg/easyredir/pagination"

type Options struct {
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
