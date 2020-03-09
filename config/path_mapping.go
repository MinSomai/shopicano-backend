package config

import (
	"github.com/spf13/viper"
)

type PathMapping struct {
	Paths map[string]string
}

var pathMapping PathMapping

func PathMappingCfg() map[string]string {
	return pathMapping.Paths
}

func LoadPathMapping() {
	mu.Lock()
	defer mu.Unlock()

	pathMapping = PathMapping{
		Paths: viper.GetStringMapString("paths_mapping"),
	}
}
