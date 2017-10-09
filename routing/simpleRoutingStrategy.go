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

func (srs SimpleRoutingStrategy) LocateClientService(sr *core.ServerRequest) (*core.ClientRequest, error) {
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

func (srs SimpleRoutingStrategy) GetServiceKind(cr *core.ClientRequest) *core.ServiceKind {
	sk := &core.ServiceKind{}
	sk.Protocol = coretype.Http
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
	countlog.Infof("event!select-version: %s", sk.Version)
	return sk
}

func (srs SimpleRoutingStrategy) SelectOneInst(sk *core.ServiceKind) (*core.ServiceInstance, error) {
	ipString, err := FindServiceKindAddr(sk)
	if err != nil {
		return nil, err
	}
	addr, err := net.ResolveTCPAddr("tcp", ipString)
	if err != nil {
		return nil, err
	}
	inst := &core.ServiceInstance{
		ServiceKind: sk,
		RemoteAddr:  addr,
	}
	return inst, nil
}

func (srs SimpleRoutingStrategy) GetRoutingDecision(inst *core.ServiceInstance) *core.RoutingDecision {
	return &core.RoutingDecision{
		ServiceInstance: inst,
		Verdict:         core.Accept,
		WaitDuration:    time.Millisecond * 50,
	}
}
