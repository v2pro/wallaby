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
	return net.Dial("tcp", "pypi.doubanio.com:80")
}

