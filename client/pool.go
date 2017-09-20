package client

import (
	"github.com/v2pro/wallaby/core"
	"net"
	"sync"
)

var pools = map[core.Qualifier]chan *net.TCPConn{}
var poolsMutex = &sync.Mutex{}

func getPool(qualifier core.Qualifier) chan *net.TCPConn {
	poolsMutex.Lock()
	defer poolsMutex.Unlock()
	pool := pools[qualifier]
	if pool == nil {
		pool = make(chan *net.TCPConn, 8)
		pools[qualifier] = pool
	}
	return pool
}