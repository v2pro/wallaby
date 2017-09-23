package core

import (
	"net"
	"github.com/v2pro/wallaby/core/codec"
	"time"
)

// Overall Proxy Sequence
// ServerConn => RoutingMode => ServerRequest => ServiceKinds => ServiceInstance => RoutingDecision
// there are three modes
// per connection routing: RoutingDecision is determined by ServerConn
//		this mode is most generic, can handle any kind of tcp stream without knowing the protocol
// per stream routing: RoutingDecision is determined by first request packet in the connection
// 		or stream (when protocol is multiplex, there might be multiple streams on one connection)
//		this mode do not need to do stateful protocol handling, and can route with more information
// per packet routing (a.k.a RPC mode): RoutingDecision might be different for different request packet
//		this mode is most powerful and most costly, need complete implementation of protocol
//		including encoding/decoding/stateful action sequences

// accept ServerConn from inbound
// the LocalAddr and RemoteAddr is from tcp connection
// src service name might be bundled with the proxy installation (serving only one service)
// src service name might be inferred from RemoteAddr or LocalAddr
// dst service name might be inferred from RemoteAddr or LocalAddr

type ServerConn struct {
	LocalAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

// ServerConn => RoutingMode decision point
// we may handle different incoming port using different stream forwarding mode

type RoutingMode string

const PerConnection RoutingMode = "PerConnection"
const PerStream RoutingMode = "PerStream"
const PerPacket RoutingMode = "PerPacket"

type ConnForwardingDecision struct {
	RoutingMode RoutingMode
	ServerProtocol string
}

// read ServerRequest from inbound
// Packet might be nil if routing mode is per connection
// src/dst service is extracted from packet, might be empty

type ServerRequest struct {
	ServerConn *ServerConn
	ConnForwardingDecision *ConnForwardingDecision
	Packet codec.Packet
}

// ServerRequest => ServiceKinds decision point
// consult naming server and cluster routing table to find out following info
// ServiceCluster: we might redirect traffic from one data center to another
// ServiceProtocol: client target might speak different protocol,
// 		we can not assume inbound protocol is same as outbound protocol

// there might be many instances for one ServiceKind, as long as the four values are the same
// the instances are interchangeable (no big performance difference, same config, same geo-location)
// one ServerRequest might have many ServiceKinds as viable options

type ServiceKind struct {
	// what service actually is, determined by its source code
	ServiceName string
	// traffic segregation by src/type/dst etc,
	// data center is the most often used clustering strategy
	// clusters are defined for management reasons
	ServiceCluster string
	// multiple versions of the source code might be running concurrently
	// to support service version roll-out and roll-back without re-deployment
	ServiceVersion string
	// one running service os process might speak more than one protocol on different tcp ports
	ServiceProtocol string
}

// ServiceKinds => RoutingDecision decision point
// should we accept/reject/wait
// if accept, which service instance (among kinds and instances) to handle it
// the instance is picked considering the most optimal ServiceKind and most optimal instance within that kind

type ServiceInstance struct {
	ServiceKind ServiceKind
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

func (srv *ServiceKind) String() string {
	return srv.ServiceName + "-" + srv.ServiceCluster + "@" + srv.ServiceVersion
}