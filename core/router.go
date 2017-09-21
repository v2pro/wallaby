package core

func RouteServerConn(serverConn *ServerConn) *Stream {
	return &Stream{}
}

func RouteServerRequest(serverRequest *ServerRequest) *RoutingDecision {
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
