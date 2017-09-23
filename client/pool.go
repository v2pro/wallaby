package client

import (
	"github.com/v2pro/wallaby/core"
	"sync"
)

var pools = map[string]chan *pooledClient{}
var poolsMutex = &sync.Mutex{}

func getPool(qualifier *core.ServiceKind) chan *pooledClient {
	poolsMutex.Lock()
	defer poolsMutex.Unlock()
	pool := pools[qualifier.String()]
	if pool == nil {
		pool = make(chan *pooledClient, 8)
		pools[qualifier.String()] = pool
	}
	return pool
}