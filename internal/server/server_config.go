package server

import (
	"os"

	"github.com/BON4/timedQ/internal/manager"
	"github.com/BON4/timedQ/pkg/ttlstore"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	AppConfig struct {
		Port    string `yaml:"port"`
		LogFile string `yaml:"log-file"`
	} `yaml:"app"`

	ManagerCfg manager.ManagerConfig   `yaml:"manager"`
	StoreCfg   ttlstore.TTLStoreConfig `yaml:"store"`
}

func LoadServerConfig(path string) (ServerConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return ServerConfig{}, err
	}
	defer f.Close()

	var cfg ServerConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return ServerConfig{}, err
	}

	return cfg, nil
}
