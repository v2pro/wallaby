package core

func HowToRoute(serverConn *ServerConn) RoutingMode {
	return ""
}

func Route(serverRequest *ServerRequest) *RoutingDecision {
	return &RoutingDecision{
		ServiceInstance: &ServiceInstance{
			Service: Service{
				ServiceName:    "default",
				ServiceCluster: "localhost",
				ServiceVersion: "default",
			},
		},
	}
}
