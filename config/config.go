package config

import (
	"bytes"
	_ "embed"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var (
	//go:embed config.yml
	configBytes []byte
)

type Remote struct {
	TLS struct {
		ClientCertFile string
		ClientKeyFile  string
	}
	// Настройки прокси.
	Proxy struct {
		// Адрес формата http(s)://proxy.example.com:1234
		Url string
	}
}

type Config struct {
	Remote Remote
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yml")
	if err := v.ReadConfig(bytes.NewBuffer(configBytes)); err != nil {
		return nil, err
	}

	decodeHooks := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)

	var c Config
	if err := v.Unmarshal(&c, viper.DecodeHook(decodeHooks)); err != nil {
		return nil, err
	}

	return &c, nil
}
