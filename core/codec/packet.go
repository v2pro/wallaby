package codec

import "io"

var Decoders = map[string]Decoder{
	"HTTP": &httpDecoder{},
}

type Packet interface {
	Feature() map[string]string
	Write(io.Writer) error
}

type Decoder interface {
	DecodeRequest(capture *Capture) (Packet, error)
	DecodeResponse(capture *Capture) (Packet, error)
}
