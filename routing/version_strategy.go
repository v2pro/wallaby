package routing

import (
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/core/codec"
	"github.com/v2pro/wallaby/core/coretype"
	"net"
	"time"
)

var LOCALHOST = "localhost"

// VersionRoutingStrategy is for http only
type VersionRoutingStrategy struct {
	versions        *ServiceVersions
	versionsHandler *InboundService
	serviceName     string
}

func NewVersionRoutingStrategy(service string, filePath string, handlerAddr string) *VersionRoutingStrategy {
	var versions = NewServiceVersions(filePath)
	if versions == nil {
		countlog.Info("event!select-version", "NewVersionRoutingStrategy ", service)
		return nil
	}
	err := versions.Start()
	if err != nil {
		countlog.Info("event!select-version", "NewVersionRoutingStrategy ", "start fail")
		return nil
	}
	var versionsHandler = NewInboundService(handlerAddr, versions)
	versionsHandler.Start()
	return &VersionRoutingStrategy{
		versions:        versions,
		serviceName:     service,
		versionsHandler: versionsHandler,
	}
}

func (vrs *VersionRoutingStrategy) ServiceVersions() *ServiceVersions {
	return vrs.versions
}

func (vrs *VersionRoutingStrategy) Route(packet codec.Packet) *ServiceVersion {
	sv := vrs.versions.Route(packet)
	if sv != nil {
		countlog.Debug("event!Route", "addr", sv.Address)
	}
	return sv
}

// LocateClientService simply get service from wallaby_services.json
func (vrs *VersionRoutingStrategy) LocateClientService(sr *core.ServerRequest) (*core.ClientRequest, error) {
	clientService := &core.ClientRequest{}
	clientService.ServerRequest = sr
	// to localhost service
	clientService.SrcServiceName = LOCALHOST
	clientService.SrcServiceCluster = LOCALHOST
	clientService.DstServiceName = vrs.serviceName
	return clientService, nil
}

// GetServiceKind decides the version of service for the given request and return a ServiceKind of that version
func (vrs *VersionRoutingStrategy) GetServiceKind(cr *core.ClientRequest) *core.ServiceKind {
	sk := &core.ServiceKind{}
	sk.Packet = cr.ServerRequest.Packet
	sk.Protocol = coretype.HTTP
	sk.Name = cr.DstServiceName
	sk.Cluster = LOCALHOST
	return sk
}

// SelectOneInst just read the ip address from  wallaby_services.json and return ServiceInstance
func (vrs *VersionRoutingStrategy) SelectOneInst(sk *core.ServiceKind) (*core.ServiceInstance, error) {
	var ver *ServiceVersion = vrs.Route(sk.Packet)

	if ver == nil {
		return nil, nil
	}

	addr, err := net.ResolveTCPAddr("tcp", ver.Address)
	if err != nil {
		return nil, err
	}
	// set the address and version back ServiceKind
	sk.Name = ver.Address
	sk.Version = ver.Version
	countlog.Debug("event!select-version", "version", addr)
	return &core.ServiceInstance{
		Kind:       sk,
		RemoteAddr: addr,
	}, nil
}

// GetRoutingDecision return the default RoutingDecision
func (vrs *VersionRoutingStrategy) GetRoutingDecision(inst *core.ServiceInstance) *core.RoutingDecision {
	return &core.RoutingDecision{
		ServiceInstance: inst,
		Verdict:         core.Accept,
		WaitDuration:    time.Millisecond * 50,
	}
}

func (vrs *VersionRoutingStrategy) Close() error {
	return vrs.versionsHandler.Shutdown()
}
