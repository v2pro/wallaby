package coretype

// Protocol is the type of protocols wallaby supports currently
type Protocol string

const (
	// HTTP protocol
	HTTP Protocol = "http"

	// THRIFT protocol
	THRIFT Protocol = "thrift"

	// RPC protocol, needs to be specified later
	RPC Protocol = "rpc"
)
