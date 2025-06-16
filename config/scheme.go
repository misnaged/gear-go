package config

// Scheme represents the application configuration scheme.
type Scheme struct {
	// Env is the application environment.
	Env    string
	Client *Client
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
