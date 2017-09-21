package client

import (
	"github.com/v2pro/wallaby/core"
	"sync"
)

var pools = map[core.Service]chan *pooledClient{}
var poolsMutex = &sync.Mutex{}

func getPool(qualifier core.Service) chan *pooledClient {
	poolsMutex.Lock()
	defer poolsMutex.Unlock()
	pool := pools[qualifier]
	if pool == nil {
		pool = make(chan *pooledClient, 8)
		pools[qualifier] = pool
	}
	return pool
}