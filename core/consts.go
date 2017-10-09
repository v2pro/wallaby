package core

type RoutingModeType string
type Verdict int

const (

	PerConnection RoutingModeType = "PerConnection"
	PerStream RoutingModeType = "PerStream"
	PerPacket RoutingModeType = "PerPacket"

	Accept Verdict = 1
	Reject Verdict = 2
	Wait Verdict = 3
)


