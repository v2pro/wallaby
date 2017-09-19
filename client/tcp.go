package client

import "net"

type tcpOutboundClient struct {
	*net.TCPConn
}

