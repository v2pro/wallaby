package client

import (
	"io"
	"github.com/v2pro/wallaby/core"
	"net"
	"go.uber.org/atomic"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/core/codec"
	"bufio"
)

type Client interface {
	Handle(req codec.Packet, capture *codec.Capture) (codec.Packet, error)
	io.Closer
}

type pooledClient struct {
	*net.TCPConn
	pool      chan *pooledClient
	isInvalid atomic.Bool
	codec     codec.Codec
	reader    *bufio.Reader
}

func (clt *pooledClient) Handle(req codec.Packet, capture *codec.Capture) (codec.Packet, error) {
	err := clt.codec.EncodeRequest(req, clt)
	if err != nil {
		clt.isInvalid.Store(true)
		countlog.Warn("event!client.failed to write request", "err", err)
		return nil, err
	}
	err = capture.SetReader(clt)
	if err != nil {
		countlog.Warn("event!client.capture has remaining buffer", "err", err)
		return nil, err
	}
	resp, err := clt.codec.DecodeResponse(capture)
	if err != nil {
		clt.isInvalid.Store(true)
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
	case clt.pool <- clt:
		countlog.Trace("event!client.pool_recycle",
			"conn", clt.TCPConn.LocalAddr())
		return nil
	default:
		countlog.Trace("event!client.pool_overflow",
			"conn", clt.TCPConn.LocalAddr())
		return clt.TCPConn.Close()
	}
}

func Get(target *core.ServiceInstance) (Client, error) {
	pool := getPool(target.ServiceKind)
	select {
	case client := <-pool:
		countlog.Trace("event!client.reuse",
			"qualifier", target.ServiceKind.String(),
			"conn", client.TCPConn.LocalAddr())
		return client, nil
	default:
		return GetNew(target.ServiceKind)
	}
}

func GetNew(qualifier *core.ServiceKind) (Client, error) {
	pool := getPool(qualifier)
	return connect(pool, qualifier)
}

func connect(pool chan *pooledClient, qualifier *core.ServiceKind) (Client, error) {
	addr, err := core.FindServiceAddr(qualifier)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	countlog.Trace("event!client.connect",
		"qualifier", qualifier.String(),
		"conn", conn.LocalAddr())
	clt := &pooledClient{
		TCPConn: conn.(*net.TCPConn),
		pool:    pool,
		codec:   codec.Codecs["HTTP"],
	}
	clt.reader = bufio.NewReaderSize(clt, 2048)
	return clt, nil
}
