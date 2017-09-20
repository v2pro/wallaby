package codec

import (
	"net/http"
	"io"
	"fmt"
	"reflect"
	"bufio"
)

type httpRequestPacket struct {
	*http.Request
}

func (req *httpRequestPacket) GetFeature(key string) string {
	return ""
}

type httpResponsePacket struct {
	*http.Response
}

func (resp *httpResponsePacket) GetFeature(key string) string {
	return ""
}

type httpCodec struct {
}

func (codec *httpCodec) DecodeRequest(reader *bufio.Reader) (Packet, error) {
	httpReq, err := http.ReadRequest(reader)
	if err != nil {
		return nil, err
	}
	return &httpRequestPacket{
		Request: httpReq,
	}, nil
}

func (codec *httpCodec) DecodeResponse(reader *bufio.Reader) (Packet, error) {
	httpResp, err := http.ReadResponse(reader, nil)
	if err != nil {
		return nil, err
	}
	return &httpResponsePacket{
		Response: httpResp,
	}, nil
}

func (codec *httpCodec) EncodeRequest(request Packet, writer io.Writer) error {
	switch typed := request.(type) {
	case *httpRequestPacket:
		return typed.Request.Write(writer)
	default:
		return fmt.Errorf("http codec can not encode request of type: " + reflect.TypeOf(request).String())
	}
}

func (codec *httpCodec) EncodeResponse(response Packet, writer io.Writer) error {
	switch typed := response.(type) {
	case *httpResponsePacket:
		return typed.Response.Write(writer)
	default:
		return fmt.Errorf("http codec can not encode response of type: " + reflect.TypeOf(response).String())
	}
}
