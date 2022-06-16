package easyredir

import (
	"net/http"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir/client"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/host"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/option"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/rule"
)

type Easyredir struct {
	Client *client.Client
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

func New(opts ...option.Option) *Easyredir {
	return &Easyredir{
		Client: client.New(opts...),
	}
}

func (c *Easyredir) CreateRule(attr rule.Attributes, opts ...option.Option) (r rule.Rule, err error) {
	return rule.CreateRule(c.Client, attr, opts...)
}

func (c *Easyredir) ListRules(opts ...option.Option) (r rule.Rules, err error) {
	return rule.ListRulesPaginator(c.Client, opts...)
}

func (c *Easyredir) RemoveRule(id string) (res bool, err error) {
	return rule.RemoveRule(c.Client, id)
}

func (c *Easyredir) GetHost(id string) (h host.Host, err error) {
	return host.GetHost(c.Client, id)
}

func (c *Easyredir) ListHosts(opts ...option.Option) (h host.Hosts, err error) {
	return host.ListHostsPaginator(c.Client, opts...)
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

type WithInclude string

func (i WithInclude) Apply(o *option.Options) {
	o.Include = string(i)
}

type WithHTTPClient struct {
	client Doer
}

func (c WithHTTPClient) Apply(o *option.Options) {
	o.HTTPClient = c.client
}
