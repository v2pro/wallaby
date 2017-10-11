package codec

import (
	"github.com/v2pro/wallaby/core/coretype"
	"io"
)

// Codecs is the codecs map
var Codecs = map[coretype.Protocol]Codec{
	coretype.HTTP: &httpCodec{},
}

// Packet defines methods for msg body
type Packet interface {
	GetFeature(key string) string
}

// Codec defines encode/decode methods for request/response
type Codec interface {

	// DecodeRequest reads from capture, and returns a Packet obj and error
	DecodeRequest(capture *Capture) (Packet, error)

	// DecodeResponse reads from capture and returns a Packet obj and error
	DecodeResponse(capture *Capture) (Packet, error)

	// EncodeRequest encodes the request and writes to destination(writer)
	EncodeRequest(request Packet, writer io.Writer) error

	// EncodeResponse encodes the response and writes to destination(writer)
	EncodeResponse(response Packet, writer io.Writer) error
}
