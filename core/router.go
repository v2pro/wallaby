package core

func HowToRoute(serverConn *ServerConn) RoutingMode {
	return ""
}

func Route(serverRequest *ServerRequest) *RoutingDecision {
	return &RoutingDecision{
		ServiceInstance: &ServiceInstance{
			ServiceKind: ServiceKind{
				ServiceName:    "default",
				ServiceCluster: "localhost",
				ServiceVersion: "default",
			},
		},
	}
}
