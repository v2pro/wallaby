package client

import (
	"io"
	"github.com/v2pro/wallaby/core"
	"net"
)


type OutboundClient interface {
	io.Writer
	io.Reader
	Close() error
}

func Connect(qualifier core.OutboundQualifier) (OutboundClient, error) {
	return net.Dial("tcp", "127.0.0.1:8849")
}

