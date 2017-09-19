package core

func Route(request InboundRequest) OutboundQualifier {
	return OutboundQualifier{
		ServiceName: "default",
		ServiceDC: "localhost",
		ServiceVersion: "default",
	}
}