package codec

import (
	"io"
)

var Codecs = map[string]Codec{
	"HTTP": &httpCodec{},
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
