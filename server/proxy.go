package server

import (
	"net"
	"github.com/v2pro/wallaby/countlog"
	"github.com/v2pro/wallaby/core/codec"
)

func Start() {
	addr := "127.0.0.1:8848"
	decoder := codec.Codecs["HTTP"]
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
		go newStream(conn.(*net.TCPConn), decoder).proxy()
	}
}
