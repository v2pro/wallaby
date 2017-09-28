package core

import (
	"net"
	"time"

	"github.com/v2pro/wallaby/core/codec"
	"github.com/v2pro/wallaby/core/coretype"
)

// ServerConn accept connection from inbound
type ServerConn struct {
	LocalAddr  *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

// ServerConn => ConnForwardingDecision: when a tcp connection is established, how to forward the connection (routing mode/protocol) is determined
type ConnForwardingDecision struct {
	RoutingMode    RoutingModeType
	ServerProtocol coretype.Protocol
}

// ConnForwardingDecision => ServerRequest: parse request arrived server
// Packet might be nil if routing mode is per connection
type ServerRequest struct {
	ServerConn             *ServerConn
	ConnForwardingDecision *ConnForwardingDecision
	Packet                 codec.Packet
}

// ServerRequest => ClientRequest: by parsing the request, we know what is the target service
type ClientRequest struct {
	ServerRequest     *ServerRequest
	SrcServiceName    string
	SrcServiceCluster string
	DstServiceName    string
}

// ClientRequest => ServiceKinds: one service have many clusters, filter out feasible clusters by cluster routing table.
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
	Protocol coretype.Protocol
}

// ServiceKinds => ServiceInstance: choose one most optimal service cluster from many clusters,
// choose one most optimal service instance from that cluster.
// there might be many instances for one ServiceKind, as long as the four values are the same
// the instances are interchangeable (no big performance difference, same config, same geo-location)
type ServiceInstance struct {
	ServiceKind *ServiceKind
	RemoteAddr  *net.TCPAddr
}

// ServiceInstance => RoutingDecision: given the service status, should we accept/reject/wait the request.
// If accept, handle the request by chosen service instance.
type RoutingDecision struct {
	ServiceInstance *ServiceInstance
	Verdict         Verdict
	RejectResponse  interface{}
	WaitDuration    time.Duration
}

func (srv *ServiceKind) String() string {
	return srv.Name + "-" + srv.Cluster + "@" + srv.Version
}
