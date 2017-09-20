package core

import "github.com/v2pro/wallaby/core/codec"

func Route(request codec.Packet) (Qualifier, codec.Decoder) {
	return Qualifier{
		ServiceName: "default",
		ServiceDC: "localhost",
		ServiceVersion: "default",
	}, codec.Decoders["HTTP"]
}