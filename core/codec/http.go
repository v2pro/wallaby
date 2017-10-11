package codec

import (
	"net/http"
	"io"
	"fmt"
	"reflect"
	"io/ioutil"
)

type httpRequestPacket struct {
	*http.Request
	body []byte
	origReq []byte
}

func (req *httpRequestPacket) GetFeature(key string) string {
	return ""
}

type httpResponsePacket struct {
	*http.Response
	body []byte
	origResp []byte
}

func (resp *httpResponsePacket) GetFeature(key string) string {
	return ""
}

type httpCodec struct {
}

// DecodeRequest reads the http request from capture, and returns a Packet obj
func (codec *httpCodec) DecodeRequest(capture *Capture) (Packet, error) {
	httpReq, err := http.ReadRequest(capture.Reader())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		return nil, err
	}
	return &httpRequestPacket{
		Request: httpReq,
		body: body,
		origReq: capture.Bytes(),
	}, nil
}

// DecodeResponse reads from capture(http body), and returns a Packet obj
func (codec *httpCodec) DecodeResponse(capture *Capture) (Packet, error) {
	httpResp, err := http.ReadResponse(capture.Reader(), nil)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	return &httpResponsePacket{
		Response: httpResp,
		body: body,
		origResp: capture.Bytes(),
	}, nil
}

// EncodeRequest encodes the request and writes to destination(writer)
func (codec *httpCodec) EncodeRequest(request Packet, writer io.Writer) error {
	switch typed := request.(type) {
	case *httpRequestPacket:
		_, err := writer.Write(typed.origReq)
		return err
	default:
		return fmt.Errorf("http codec can not encode request of type: " + reflect.TypeOf(request).String())
	}
}

// EncodeResponse encodes the response and writes to destination(writer)
func (codec *httpCodec) EncodeResponse(response Packet, writer io.Writer) error {
	switch typed := response.(type) {
	case *httpResponsePacket:
		_, err := writer.Write(typed.origResp)
		return err
	default:
		return fmt.Errorf("http codec can not encode response of type: " + reflect.TypeOf(response).String())
	}
}
