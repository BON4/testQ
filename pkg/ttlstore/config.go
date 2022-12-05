package ttlstore

import (
	"time"
)

type TTLStoreConfig struct {
	GCRefresh time.Duration `yaml:"gc-refresh-time" mapstructure:"GC_REFRESH_TIME"`
	GCWorkers uint          `yaml:"gc-workers-num" mapstructure:"GC_WORKERS_NUM"`
	SavePath  string        `yaml:"save-path" mapstructure:"SAVE_PATH"`
	Save      bool          `yaml:"save" mapstructure:"SAVE"`
}

func NewMapStoreConfig(GCRefresh time.Duration, GCWorkers uint, path string, save bool) TTLStoreConfig {

	return TTLStoreConfig{
		GCRefresh: time.Second / 3,
		GCWorkers: 1,
		SavePath:  path,
		Save:      save,
	}
}
