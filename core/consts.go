package core

type ProtocolType string
type RoutingMode string
type Verdict int

const (
	Http ProtocolType = "http"
	Thrift ProtocolType = "thrift"
	Rpc ProtocolType = "rpc"

	PerConnection RoutingMode = "PerConnection"
	PerStream RoutingMode = "PerStream"
	PerPacket RoutingMode = "PerPacket"

	Accept Verdict = 1
	Reject Verdict = 2
	Wait Verdict = 3
)


