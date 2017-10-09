package server

import (
	"github.com/v2pro/wallaby/config"
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/core/codec"
	"github.com/v2pro/plz/countlog"
	"net"
)

func Start() {
	addr := config.ProxyAddr
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
		serverConn := &core.ServerConn{
			LocalAddr:  conn.LocalAddr().(*net.TCPAddr),
			RemoteAddr: conn.RemoteAddr().(*net.TCPAddr),
		}
		connForwardingDecision := core.HowToForward(serverConn)
		switch connForwardingDecision.RoutingMode {
		case core.PerPacket:
			decoder := codec.Codecs[connForwardingDecision.ServerProtocol]
			go newStream(conn.(*net.TCPConn), decoder).proxy()
		default:
			panic("RoutingMode not supported yet: " + connForwardingDecision.RoutingMode)
		}
	}
}
