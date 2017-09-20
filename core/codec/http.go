package codec

import (
	"net/http"
	"io"
)

type httpRequestPacket struct {
	*http.Request
	origReq []byte
}

func (req *httpRequestPacket) Feature() map[string]string {
	return map[string]string{}
}

func (req *httpRequestPacket) Write(writer io.Writer) error {
	_, err := writer.Write(req.origReq)
	return err
}

type httpResponsePacket struct {
	*http.Response
	origResp []byte
}

func (resp *httpResponsePacket) Feature() map[string]string {
	return map[string]string{}
}

func (resp *httpResponsePacket) Write(writer io.Writer) error {
	_, err := writer.Write(resp.origResp)
	return err
}

type httpDecoder struct {
}

func (decoder *httpDecoder) DecodeRequest(capture *Capture) (Packet, error) {
	httpReq, err := http.ReadRequest(capture.Reader())
	if err != nil {
		return nil, err
	}
	return &httpRequestPacket{
		Request: httpReq,
		origReq: capture.Bytes(),
	}, nil
}

func (decoder *httpDecoder) DecodeResponse(capture *Capture) (Packet, error) {
	httpResp, err := http.ReadResponse(capture.Reader(), nil)
	if err != nil {
		return nil, err
	}

	return &httpResponsePacket{
		Response: httpResp,
		origResp: capture.Bytes(),
	}, nil
}
