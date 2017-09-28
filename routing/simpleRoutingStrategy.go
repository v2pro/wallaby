package routing

import (
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/core/coretype"

	"net"
	"strconv"
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
	// we can choose the version based on deviceID/ip
	//deviceID := cr.ServerRequest.Packet.GetFeature("device-id")
	ipFrom := cr.ServerRequest.Packet.GetFeature("x-forwarded-for")
	digit := "0"
	if len(ipFrom) > 2 {
		digit = ipFrom[len(ipFrom)-2:]
	}
	percent, _ := strconv.Atoi(digit)
	if percent >= 50 {
		sk.Version = GetNextVersion()
	} else {
		sk.Version = GetCurrentVersion()
	}
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
