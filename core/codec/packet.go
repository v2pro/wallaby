package codec

import (
	"github.com/v2pro/wallaby/core/coretype"
	"io"
)

var Codecs = map[coretype.Protocol]Codec{
	coretype.Http: &httpCodec{},
}

type Packet interface {
	GetFeature(key string) string
}

type Codec interface {
	DecodeRequest(capture *Capture) (Packet, error)
	DecodeResponse(capture *Capture) (Packet, error)
	EncodeRequest(request Packet, writer io.Writer) error
	EncodeResponse(response Packet, writer io.Writer) error
}
