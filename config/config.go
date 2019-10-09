package config

import (
	"errors"
	"fmt"
	"sync"

	"github.com/spf13/viper"
	// this package is necessary to read config from remote consul
	_ "github.com/spf13/viper/remote"
)

var mu sync.Mutex

// LoadConfig initiates of config load
func LoadConfig() error {
	if err := viper.BindEnv("CONSUL_URL"); err != nil {
		return err
	}
	if err := viper.BindEnv("CONSUL_PATH"); err != nil {
		return err
	}

	consulURL := viper.GetString("CONSUL_URL")
	consulPath := viper.GetString("CONSUL_PATH")
	if consulURL == "" {
		return errors.New("CONSUL_URL missing")
	}
	if consulPath == "" {
		return errors.New("CONSUL_PATH missing")
	}

	if err := viper.AddRemoteProvider("consul", consulURL, consulPath); err != nil {
		return err
	}
	viper.SetConfigType("yml")

	if err := viper.ReadRemoteConfig(); err != nil {
		return errors.New(fmt.Sprintf("%s named \"%s\"", err.Error(), consulPath))
	}

	LoadApp()
	LoadDB()
	LoadMinio()
	LoadPaymentGateway()

	return nil
}
