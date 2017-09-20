package core

type OutboundQualifier struct {
	ServiceName string
	ServiceDC string
	ServiceVersion string
}

func (qualifier *OutboundQualifier) String() string {
	return qualifier.ServiceName + "-" + qualifier.ServiceDC + "@" + qualifier.ServiceVersion
}