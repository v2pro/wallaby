package client

import (
	"github.com/v2pro/wallaby/config"
	"github.com/v2pro/wallaby/core"
	"sync"
)

var pools = map[string]chan *pooledClient{}
var poolsMutex = &sync.Mutex{}

func getPool(sk *core.ServiceKind) chan *pooledClient {
	poolsMutex.Lock()
	defer poolsMutex.Unlock()
	pool := pools[sk.Qualifier()]
	if pool == nil {
		pool = make(chan *pooledClient, config.MaxClientPoolSize)
		pools[sk.Qualifier()] = pool
	}
	return pool
}
