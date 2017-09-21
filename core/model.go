package core

import (
	"net"
	"github.com/v2pro/wallaby/core/codec"
	"time"
)

// Overall Proxy Sequence
// ServerConn => Stream => ServerRequest => Service => ClientTarget
// there are three modes
// per connection routing: ServiceInstance is determined by ServerConn
// first packet routing: ServiceInstance is determined by first request packet
// per packet routing: ServiceInstance might be different for different request packet
// notice: the "connection" here is logical, it might not be tcp connection
// in gRPC or other multiplex protocol, one tcp connection can serve multiple logic connections

// accept ServerConn from inbound
// the LocalAddr and RemoteAddr is from tcp connection
// src service name might be bundled with the proxy installation (serving only one service)
// src service name might be inferred from RemoteAddr or LocalAddr
// dst service name might be inferred from RemoteAddr or LocalAddr

type ServerConn struct {
	SrcServiceName string
	SrcServiceCluster string
	DstServiceName string
	LocalAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

// ServerConn => Stream routing decision point
// we may handle different incoming port using different stream forwarding mode
// some protocol will bind to specific port

type Stream struct {
	StreamMode string
	ServerProtocol string
}

// read ServerRequest from inbound
// Packet might be nil if not in RPC mode
// src/dst service is extracted from packet, might be empty

type ServerRequest struct {
	ServerConn *ServerConn
	Stream *Stream
	SrcServiceName string
	SrcServiceCluster string
	DstServiceName string
	Packet codec.Packet
}

// ServerRequest => Service routing decision point
// consult naming server and cluster routing table to find out following info
// ServiceCluster: we might redirect traffic from one data center to another
// ServiceProtocol: client target might speak different protocol,
// 	we can not assume inbound protocol is same as outbound protocol

// there might be many instances for one service, as long as the four values are the same
// the instances are interchangeable

type Service struct {
	ServiceName string
	ServiceCluster string
	ServiceVersion string
	ServiceProtocol string
}

// Service => RoutingDecision routing decision point
// should we accept/reject/wait
// if accept, which service instance (among clusters and instances) to handle it

type ServiceInstance struct {
	Service Service
	RemoteAddr *net.TCPAddr
}

type Verdict int

const Accept Verdict = 1
const Reject Verdict = 2
const Wait Verdict = 3

type RoutingDecision struct {
	ServiceInstance *ServiceInstance
	Verdict Verdict
	RejectResponse interface{}
	WaitDuration time.Duration
}

func (srv *Service) String() string {
	return srv.ServiceName + "-" + srv.ServiceCluster + "@" + srv.ServiceVersion
}