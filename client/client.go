package client

import (
	"io"
	"github.com/v2pro/wallaby/core"
	"net"
	"go.uber.org/atomic"
	"github.com/v2pro/wallaby/countlog"
	"github.com/v2pro/wallaby/core/codec"
)

type Client interface {
	Handle(req codec.Packet, capture *codec.Capture, decoder codec.Decoder) (codec.Packet, error)
	io.Closer
}

type pooledClient struct {
	*net.TCPConn
	pool      chan *net.TCPConn
	isInvalid atomic.Bool
}

func (clt *pooledClient) Handle(req codec.Packet, capture *codec.Capture, decoder codec.Decoder) (codec.Packet, error) {
	err := req.Write(clt)
	if err != nil {
		countlog.Warn("event!client.failed to write request", "err", err)
		return nil, err
	}
	capture.Reset(clt)
	resp, err := decoder.DecodeResponse(capture)
	if err != nil {
		countlog.Warn("event!server.failed to read response", "err", err)
		return nil, err
	}
	return resp, nil
}

func (clt *pooledClient) Read(p []byte) (n int, err error) {
	n, err = clt.TCPConn.Read(p)
	if err != nil {
		clt.isInvalid.Store(true)
	}
	return
}

func (clt *pooledClient) Write(p []byte) (n int, err error) {
	n, err = clt.TCPConn.Write(p)
	if err != nil {
		clt.isInvalid.Store(true)
	}
	return
}

func (clt *pooledClient) Close() error {
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

func Get(qualifier core.Qualifier) (Client, error) {
	pool := getPool(qualifier)
	select {
	case conn := <-pool:
		countlog.Trace("event!client.reuse",
			"qualifier", &qualifier,
			"conn", conn.LocalAddr())
		return &pooledClient{
			TCPConn: conn,
			pool:    pool,
		}, nil
	default:
		return GetNew(qualifier)
	}
}

func GetNew(qualifier core.Qualifier) (Client, error) {
	pool := getPool(qualifier)
	return connect(pool, qualifier)
}

func connect(pool chan *net.TCPConn, qualifier core.Qualifier) (Client, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8849")
	if err != nil {
		return nil, err
	}
	countlog.Trace("event!client.connect",
		"qualifier", &qualifier,
		"conn", conn.LocalAddr())
	return &pooledClient{
		TCPConn: conn.(*net.TCPConn),
		pool:    pool,
	}, nil
}
