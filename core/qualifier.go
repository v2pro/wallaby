package core

type Qualifier struct {
	ServiceName string
	ServiceDC string
	ServiceVersion string
}

func (qualifier *Qualifier) String() string {
	return qualifier.ServiceName + "-" + qualifier.ServiceDC + "@" + qualifier.ServiceVersion
}