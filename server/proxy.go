package server

import (
	"net"
	"github.com/v2pro/wallaby/countlog"
)

var decoders = map[string]decoder{
	"http": &httpDecoder{},
}

func Start() {
	addr := "127.0.0.1:8848"
	decoder := decoders["http"]
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		countlog.Error("event!server.failed to bind proxy port", "err", err)
		return
	}
	countlog.Info("event!server.started", "addr", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			countlog.Error("event!server.failed to accept outbound", "err", err)
			return
		}
		srm := &stream{
			svr: conn.(*net.TCPConn),
			decoder: decoder,
		}
		go srm.proxy()
	}
}
