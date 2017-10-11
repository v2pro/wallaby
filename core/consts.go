package core

// RoutingModeType is the methods of forwarding the connection
type RoutingModeType string

// Verdict accept/reject/wait
type Verdict int

const (
	// PerConnection not supported yet
	PerConnection RoutingModeType = "PerConnection"

	// PerStream not supported yet
	PerStream RoutingModeType = "PerStream"

	// PerPacket the most common type, like http
	PerPacket RoutingModeType = "PerPacket"

	// Accept ready to receive requests
	Accept Verdict = 1

	// Reject can not receive requests
	Reject Verdict = 2

	// Wait too busy to receive requests, have to wait
	Wait Verdict = 3
)
