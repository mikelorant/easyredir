package config

type Config struct {
	baseURL   string
	apiKey    string
	apiSecret string
}

const (
	_BaseURL = "https://api.easyredir.com/v1"
)

func New(key, secret string) *Config {
	return &Config{
		apiKey: key,
		apiSecret: secret,
		baseURL: _BaseURL,
	}
}

func (c Config) BaseURL() string {
	return c.baseURL
}

func (c Config) APIKey() string {
	return c.apiKey
}

func (c Config) APISecret() string {
	return c.apiSecret
}
