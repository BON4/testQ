package ttlstore

import (
	"time"
)

type TTLStoreConfig struct {
	GCRefresh time.Duration `yaml:"gc-refresh-time"`
	GCWorkers uint          `yaml:"gc-workers-num"`
	SavePath  string        `yaml:"save-path"`
	Save      bool          `yaml:"save"`
}

func NewMapStoreConfig(GCRefresh time.Duration, GCWorkers uint, path string, save bool) TTLStoreConfig {

	return TTLStoreConfig{
		GCRefresh: time.Second / 3,
		GCWorkers: 1,
		SavePath:  path,
		Save:      save,
	}
}
