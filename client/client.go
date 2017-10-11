package client

import (
	"bufio"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/core/codec"
	"go.uber.org/atomic"
	"io"
	"net"
)

// Client interface
type Client interface {

	// Handle req, srm.cltCapture
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

// Handle copies request to client connection, get the response from client and converts to packet and return
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

// Get fetch a Client connection
func Get(target *core.ServiceInstance) (Client, error) {
	pool := getPool(target.Kind)
	select {
	case client := <-pool:
		countlog.Trace("event!client.reuse",
			"qualifier", target.Kind.Qualifier(),
			"conn", client.TCPConn.LocalAddr())
		return client, nil
	default:
		return GetNew(target)
	}
}

// GetNew create a new Client connection
func GetNew(target *core.ServiceInstance) (Client, error) {
	pool := getPool(target.Kind)
	return connect(pool, target)
}

// connect establish a connection to a client host address and return the Client object
func connect(pool chan *pooledClient, target *core.ServiceInstance) (Client, error) {
	addr := target.RemoteAddr.String()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	countlog.Trace("event!client.connect",
		"qualifier", target.Kind.Qualifier(),
		"conn", conn.LocalAddr())
	clt := &pooledClient{
		TCPConn: conn.(*net.TCPConn),
		pool:    pool,
		codec:   codec.Codecs[target.Kind.Protocol],
	}
	clt.reader = bufio.NewReaderSize(clt, 2048)
	return clt, nil
}
