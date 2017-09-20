package client

import (
	"io"
	"github.com/v2pro/wallaby/core"
	"net"
	"go.uber.org/atomic"
	"github.com/v2pro/wallaby/countlog"
)

type OutboundClient interface {
	io.Writer
	io.Reader
	io.Closer
}

type pooledOutboundClient struct {
	*net.TCPConn
	pool      chan *net.TCPConn
	isInvalid atomic.Bool
}

func (clt *pooledOutboundClient) Read(p []byte) (n int, err error) {
	n, err = clt.TCPConn.Read(p)
	if err != nil {
		clt.isInvalid.Store(true)
	}
	return
}

func (clt *pooledOutboundClient) Write(p []byte) (n int, err error) {
	n, err = clt.TCPConn.Write(p)
	if err != nil {
		clt.isInvalid.Store(true)
	}
	return
}

func (clt *pooledOutboundClient) Close() error {
	if clt.isInvalid.Load() {
		countlog.Trace("event!client.drop_invalid",
			"conn", clt.TCPConn.LocalAddr())
		return clt.TCPConn.Close()
	}
	select {
	case clt.pool <- clt.TCPConn:
		countlog.Trace("event!client.pool_recycle",
			"conn", clt.TCPConn.LocalAddr())
		return nil
	default:
		countlog.Trace("event!client.pool_overflow",
			"conn", clt.TCPConn.LocalAddr())
		return clt.TCPConn.Close()
	}
}

func Get(qualifier core.OutboundQualifier) (OutboundClient, error) {
	pool := getPool(qualifier)
	select {
	case conn := <-pool:
		countlog.Trace("event!client.reuse",
			"qualifier", &qualifier,
			"conn", conn.LocalAddr())
		return &pooledOutboundClient{
			TCPConn: conn,
			pool:    pool,
		}, nil
	default:
		return GetNew(qualifier)
	}
}

func GetNew(qualifier core.OutboundQualifier) (OutboundClient, error) {
	pool := getPool(qualifier)
	return connect(pool, qualifier)
}

func connect(pool chan *net.TCPConn, qualifier core.OutboundQualifier) (OutboundClient, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8849")
	if err != nil {
		return nil, err
	}
	countlog.Trace("event!client.connect",
		"qualifier", &qualifier,
		"conn", conn.LocalAddr())
	return &pooledOutboundClient{
		TCPConn: conn.(*net.TCPConn),
		pool:    pool,
	}, nil
}
