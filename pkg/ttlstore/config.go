package ttlstore

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type TTLStoreConfig struct {
	MapStore struct {
		GCRefresh time.Duration `yaml:"gc-refresh-time"`
		GCWorkers uint          `yaml:"gc-workers-num"`
	} `yaml:"map"`
	RedisStore struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
}

func newMapStoreConfig(GCRefresh time.Duration, GCWorkers uint) TTLStoreConfig {

	return TTLStoreConfig{
		MapStore: struct {
			GCRefresh time.Duration `yaml:"gc-refresh-time"`
			GCWorkers uint          `yaml:"gc-workers-num"`
		}{
			GCRefresh: time.Second / 3,
			GCWorkers: 1,
		},
	}
}

func LoadConfig(path string) (TTLStoreConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return TTLStoreConfig{}, err
	}
	defer f.Close()

	var cfg TTLStoreConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return TTLStoreConfig{}, err
	}

	return cfg, nil
}
