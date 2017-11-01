package routing

import (
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/core/coretype"

	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/datacenter"
	"net"
	"time"
)

// SimpleRoutingStrategy is for http only
type SimpleRoutingStrategy struct {
}

// LocateClientService simply get service from wallaby_services.json
func (srs *SimpleRoutingStrategy) LocateClientService(sr *core.ServerRequest) (*core.ClientRequest, error) {
	clientService := &core.ClientRequest{}
	clientService.ServerRequest = sr
	serviceFrom := sr.Packet.GetFeature("x-wallaby-downstream-service")
	clusterFrom := sr.Packet.GetFeature("x-wallaby-downstream-cluster")
	clientService.SrcServiceName = serviceFrom
	clientService.SrcServiceCluster = clusterFrom
	// just pick the first one, assuming all the services are same
	node := GetCurrentServiceNode()
	clientService.DstServiceName = node.Service
	return clientService, nil
}

// GetServiceKind decides the version of service for the given request and return a ServiceKind of that version
func (srs *SimpleRoutingStrategy) GetServiceKind(cr *core.ClientRequest) *core.ServiceKind {
	sk := &core.ServiceKind{}
	sk.Protocol = coretype.HTTP
	sk.Name = cr.DstServiceName
	// try the same cluster as the upstream first
	sk.Cluster = cr.SrcServiceCluster
	if sk.Cluster == "" {
		node := GetCurrentServiceNode()
		sk.Cluster = node.Cluster
	}
	// we can choose the version based on deviceID/ip/Cityid/USN(module identifier)/...
	routingSetting := datacenter.GetRoutingSetting()
	if routingSetting.IsValid() {
		// for example, x-forwarded-for regex [12345]$, Cityid >= 10000, etc.
		hashVal := cr.ServerRequest.Packet.GetFeature(routingSetting.Hashkey)
		if routingSetting.RunRoutingRule(hashVal) {
			sk.Version = GetNextVersion()
		} else {
			sk.Version = GetCurrentVersion()
		}
	} else {
		sk.Version = GetCurrentVersion()
	}
	countlog.Info("event!select-version", "version", sk.Version)
	return sk
}

// SelectOneInst just read the ip address from  wallaby_services.json and return ServiceInstance
func (srs *SimpleRoutingStrategy) SelectOneInst(sk *core.ServiceKind) (*core.ServiceInstance, error) {
	ipString, err := FindServiceKindAddr(sk)
	if err != nil {
		return nil, err
	}
	addr, err := net.ResolveTCPAddr("tcp", ipString)
	if err != nil {
		return nil, err
	}
	inst := &core.ServiceInstance{
		Kind:       sk,
		RemoteAddr: addr,
	}
	return inst, nil
}

// GetRoutingDecision return the default RoutingDecision
func (srs *SimpleRoutingStrategy) GetRoutingDecision(inst *core.ServiceInstance) *core.RoutingDecision {
	return &core.RoutingDecision{
		ServiceInstance: inst,
		Verdict:         core.Accept,
		WaitDuration:    time.Millisecond * 50,
	}
}

func (srs *SimpleRoutingStrategy) Close() error {
	return nil
}
