package core

func Route(request Packet) OutboundQualifier {
	return OutboundQualifier{
		ServiceName: "default",
		ServiceDC: "localhost",
		ServiceVersion: "default",
	}
}