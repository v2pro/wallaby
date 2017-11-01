package core

import "github.com/v2pro/wallaby/core/coretype"

// HowToForward determines how to forward the connection (routing mode/protocol) when a tcp connection is established,
func HowToForward(serverConn *ServerConn) *ConnForwardingDecision {
	return &ConnForwardingDecision{
		RoutingMode:    PerPacket,
		ServerProtocol: coretype.HTTP,
	}
}

// HowToRoute according to the request, decide:
//      which is the service we want
//      which cluster of the service do we choose, which version of the service-cluster do we choose (ServiceKind)
//      which instance in the service-cluster do we choose
//      what is the corresponding ip:port for the chosen instance
//      how to handle the request? shall we accept, reject or wait?
// The whole process: ServerRequest => ClientRequest => ServiceKind => ServiceInstance => RoutingDecision
func HowToRoute(serverRequest *ServerRequest, rs RoutingStrategy) (*RoutingDecision, error) {
	clientService, err := rs.LocateClientService(serverRequest)
	if err != nil {
		return nil, err
	}
	sk := rs.GetServiceKind(clientService)
	inst, err := rs.SelectOneInst(sk)
	if err != nil {
		return nil, err
	}
	return rs.GetRoutingDecision(inst), nil
}

// RoutingStrategy defines the process between the core models
type RoutingStrategy interface {
	LocateClientService(sr *ServerRequest) (*ClientRequest, error)
	GetServiceKind(cr *ClientRequest) *ServiceKind
	SelectOneInst(sk *ServiceKind) (*ServiceInstance, error)
	GetRoutingDecision(inst *ServiceInstance) *RoutingDecision
	Close() error
}
