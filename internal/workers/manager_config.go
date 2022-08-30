package workers

import (
	"time"
)

type ManagerConfig struct {
	Manager struct {
		WorkerNum uint          `yaml:"worker-num"`
		ValTTL    time.Duration `ymal:"val-ttl"`
	} `yaml:"manager"`
}

func newManagerConfig(WorkerNum uint, ValTTL time.Duration) ManagerConfig {

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
