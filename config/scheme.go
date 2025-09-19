package config

import "github.com/misnaged/substrate-api-rpc/keyring"

// Scheme represents the application configuration scheme.
type Scheme struct {
	// Env is the application environment.
	Env           string
	Client        *Client
	Keyring       *Keyring
	Subscriptions *Subscriptions
}
type Subscriptions struct {
	BuffSize int
	Enabled  bool
}
type Keyring struct {
	Category keyring.Category
	Seed     string
}
type Client struct {
	Addr        `mapstructure:",squash"`
	IsWebSocket bool
	IsSecured   bool
}

type Addr struct {
	Transport string
	Host      string
	Port      int
}
