package manager

import (
	"time"
)

type ManagerConfig struct {
	WorkerNum uint          `yaml:"worker-num"`
	ValTTL    time.Duration `yaml:"val-ttl"`
}

func newManagerConfig(WorkerNum uint, ValTTL time.Duration) ManagerConfig {
	return ManagerConfig{
		ValTTL:    ValTTL,
		WorkerNum: WorkerNum,
	}

}
