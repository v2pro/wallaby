package routing

import (
	"github.com/v2pro/wallaby/core"
	"net"
)

type ServiceDiscover interface {
	FindServiceKindAddr(sk *core.ServiceKind) (*net.TCPAddr, error)
}
