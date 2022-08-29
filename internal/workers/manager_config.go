package workers

import (
	"time"

	"github.com/BON4/timedQ/pkg/ttlstore"
)

type ManagerConfig struct {
	Manager struct {
		WorkerNum uint          `yaml:"worker-num"`
		ValTTL    time.Duration `ymal:"val-ttl"`
	} `yaml:"manager"`
}

func newManagerConfig(WorkerNum uint, ValTTL time.Duration, StoreCfg ttlstore.TTLStoreConfig) ManagerConfig {

	return ManagerConfig{
		Manager: struct {
			WorkerNum uint          `yaml:"worker-num"`
			ValTTL    time.Duration `ymal:"val-ttl"`
		}{
			ValTTL:    ValTTL,
			WorkerNum: WorkerNum,
		},
	}

}
