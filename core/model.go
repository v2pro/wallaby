package core

import (
	"net"
	"time"

	"github.com/v2pro/wallaby/core/codec"
)

// accept ServerConn from inbound
// the LocalAddr and RemoteAddr is from tcp connection
// src service name might be bundled with the proxy installation (serving only one service)
// src service name might be inferred from RemoteAddr or LocalAddr
// dst service name might be inferred from RemoteAddr or LocalAddr

type ServerConn struct {
	LocalAddr  *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

// ServerConn => RoutingMode decision point
// we may handle different incoming port using different stream forwarding mode

type ConnForwardingDecision struct {
	RoutingMode    RoutingMode
	ServerProtocol string
}

// read ServerRequest from inbound
// Packet might be nil if routing mode is per connection
// src/dst service is extracted from packet, might be empty

type ServerRequest struct {
	ServerConn             *ServerConn
	ConnForwardingDecision *ConnForwardingDecision
	Packet                 codec.Packet
}

type ClientRequest struct {
	ServerRequest *ServerRequest
	SrcServiceName string
	SrcServiceCluster string
	DstServiceName string
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
	Name string
	// traffic segregation by src/type/dst etc,
	// data center is the most often used clustering strategy
	// clusters are defined for management reasons
	Cluster string
	// multiple versions of the source code might be running concurrently
	// to support service version roll-out and roll-back without re-deployment
	Version string
	// one running service os process might speak more than one protocol on different tcp ports
	Protocol ProtocolType
}

// ServiceKinds => RoutingDecision decision point
// should we accept/reject/wait
// if accept, which service instance (among kinds and instances) to handle it
// the instance is picked considering the most optimal ServiceKind and most optimal instance within that kind

type ServiceInstance struct {
	ServiceKind *ServiceKind
	RemoteAddr  *net.TCPAddr
}

type RoutingDecision struct {
	ServiceInstance *ServiceInstance
	Verdict         Verdict
	RejectResponse  interface{}
	WaitDuration    time.Duration
}

func (srv *ServiceKind) String() string {
	return srv.Name + "-" + srv.Cluster + "@" + srv.Version
}

type ServiceNode struct {
	Qualifier string `json:"qualifier"`
	Address   string `json:"address"`
}
