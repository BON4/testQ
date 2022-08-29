package config

import (
	"github.com/BON4/timedQ/internal/workers"
	"github.com/BON4/timedQ/pkg/ttlstore"
)

type ServiceConfig struct {
	ManagerCfg workers.ManagerConfig
	StoreCfg   ttlstore.TTLStoreConfig
}
