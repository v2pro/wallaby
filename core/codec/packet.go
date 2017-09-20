package codec

import (
	"io"
	"bufio"
)

var Codecs = map[string]Codec{
	"HTTP": &httpCodec{},
}

type Packet interface {
	GetFeature(key string) string
}

type Codec interface {
	DecodeRequest(reader *bufio.Reader) (Packet, error)
	DecodeResponse(reader *bufio.Reader) (Packet, error)
	EncodeRequest(request Packet, writer io.Writer) error
	EncodeResponse(response Packet, writer io.Writer) error
}
