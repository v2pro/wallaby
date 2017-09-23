package client

import (
	"github.com/v2pro/wallaby/core"
	"sync"
	"github.com/v2pro/wallaby/config"
)

var pools = map[string]chan *pooledClient{}
var poolsMutex = &sync.Mutex{}

func getPool(qualifier *core.ServiceKind) chan *pooledClient {
	poolsMutex.Lock()
	defer poolsMutex.Unlock()
	pool := pools[qualifier.String()]
	if pool == nil {
		pool = make(chan *pooledClient, config.MaxClientPoolSize)
		pools[qualifier.String()] = pool
	}
	return pool
}