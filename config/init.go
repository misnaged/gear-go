package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func InitConfig(cfg *Scheme) error {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	viper.SetDefault("env", "prod")
	viper.SetDefault("client.host", "127.0.0.1")
	viper.SetDefault("client.port", 9944)
	viper.SetDefault("client.iswebsocket", true)
	viper.SetDefault("client.issecured", false)
	viper.SetDefault("keyring.category", "Sr25519")
	viper.SetDefault("keyring.seed", "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a") //Alice
	viper.SetDefault("subscriptions.enabled", true)
	viper.SetDefault("subscriptions.buffsize", 100)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("read config file: %w", err)
		}
	}

	setTransport()

	return viper.Unmarshal(cfg)
}

func addSecured(transport string) string {
	return fmt.Sprintf("%ss", transport)
}
func setTransport() {
	var transportType string
	if viper.GetBool("client.iswebsocket") {
		transportType = "ws"
	} else {
		transportType = "http"
	}
	if viper.GetBool("client.issecured") {
		viper.Set("client.transport", addSecured(viper.GetString("client.transport")))
	}
	viper.SetDefault("client.transport", transportType)
}
