package easyredir

import (
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/client"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/config"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/host"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/rule"
)

type Easyredir struct {
	Client *client.Client
	Config *config.Config
}

func New(apiKey, apiSecret string) *Easyredir {
	cfg := config.New(apiKey, apiSecret)

	return &Easyredir{
		Client: client.New(cfg),
		Config: cfg,
	}
}

func (c *Easyredir) Ping() string {
	return "pong"
}

func (c *Easyredir) ListRules() (r rule.Rules, err error) {
	return rule.ListRulesPaginator(c.Client, c.Config)
}

func (c *Easyredir) ListHosts() (h host.Hosts, err error) {
	return host.ListHostsPaginator(c.Client, c.Config)
}

func (c *Easyredir) GetHost(id string) (h host.Host, err error) {
	return host.GetHost(c.Client, c.Config, id)
}
