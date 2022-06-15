package easyredir

import (
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/host"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/option"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/rule"
)

type Easyredir struct {
	Client *Client
}

func New(opts ...option.Option) *Easyredir {
	o := &option.Options{}

	for _, opt := range opts {
		opt.Apply(o)
	}

	return &Easyredir{
		Client: NewClient(
			WithAPIKey(o.APIKey),
			WithAPISecret(o.APISecret),
		),
	}
}

func (c *Easyredir) ListRules(opts ...option.Option) (r rule.Rules, err error) {
	return rule.ListRulesPaginator(c.Client, opts...)
}

func (c *Easyredir) ListHosts(opts ...option.Option) (h host.Hosts, err error) {
	return host.ListHostsPaginator(c.Client, opts...)
}

func (c *Easyredir) GetHost(id string) (h host.Host, err error) {
	return host.GetHost(c.Client, id)
}

type WithLimit int

func (l WithLimit) Apply(o *option.Options) {
	o.Limit = int(l)
}

type WithSourceFilter string

func (s WithSourceFilter) Apply(o *option.Options) {
	o.SourceFilter = string(s)
}

type WithTargetFilter string

func (t WithTargetFilter) Apply(o *option.Options) {
	o.TargetFilter = string(t)
}

type WithAPIKey string

func (k WithAPIKey) Apply(o *option.Options) {
	o.APIKey = string(k)
}

type WithAPISecret string

func (s WithAPISecret) Apply(o *option.Options) {
	o.APISecret = string(s)
}

type WithBaseURL string

func (u WithBaseURL) Apply(o *option.Options) {
	o.BaseURL = string(u)
}

type WithHTTPClient struct {
	client Doer
}

func (c WithHTTPClient) Apply(o *option.Options) {
	o.HTTPClient = c.client
}
