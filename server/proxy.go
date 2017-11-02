package server

import (
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/config"
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/core/codec"
	"github.com/v2pro/wallaby/routing"
	"net"
)

type ProxyServer struct {
	conn            net.Conn
	routingStrategy core.RoutingStrategy
}

// Start runs the main wallaby server, handle incoming requests and dispatch to clients
func (p *ProxyServer) Start() error {
	addr := config.ProxyAddr
	p.routingStrategy = routing.NewVersionRoutingStrategy(
		config.ProxyServiceName,
		config.ProxyServiceVersionConfig,
		config.VersionHandlerAddr,
		config.ProxyBuildTimestamp)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		countlog.Error("event!server.failed to bind proxy port", "err", err)
		return err
	}
	countlog.Info("event!server.started", "addr", addr)
	for {
		p.conn, err = listener.Accept()
		if err != nil {
			countlog.Error("event!server.failed to accept outbound", "err", err)
			return err
		}
		serverConn := &core.ServerConn{
			LocalAddr:  p.conn.LocalAddr().(*net.TCPAddr),
			RemoteAddr: p.conn.RemoteAddr().(*net.TCPAddr),
		}
		connForwardingDecision := core.HowToForward(serverConn)
		switch connForwardingDecision.RoutingMode {
		case core.PerPacket:
			decoder := codec.Codecs[connForwardingDecision.ServerProtocol]
			go newStream(p.conn.(*net.TCPConn), decoder, p.routingStrategy).proxy()
		default:
			panic("RoutingMode not supported yet: " + connForwardingDecision.RoutingMode)
		}
	}
	return nil
}

func (p *ProxyServer) Stop() {
	if p.conn != nil {
		p.conn.Close()
	}
	if p.routingStrategy != nil {
		p.routingStrategy.Close()
	}
}
