package server

import (
	"net/http"
	"bufio"
	"github.com/v2pro/wallaby/core"
)

type httpRequestPacket struct {
	*http.Request
}

func (req *httpRequestPacket) Feature() map[string]string {
	return map[string]string{}
}

type httpResponsePacket struct {
	*http.Response
}

func (resp *httpResponsePacket) Feature() map[string]string {
	return map[string]string{}
}

type httpDecoder struct {
}

func (decoder *httpDecoder) decodeRequest(reader *bufio.Reader) (core.Packet, error) {
	httpReq, err := http.ReadRequest(reader)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Connection", "keep-alive")
	httpReq.Header.Set("Keep-Alive", "timeout=4")
	httpReq.Close = false
	return &httpRequestPacket{httpReq}, nil
}

func (decoder *httpDecoder) decodeResponse(reader *bufio.Reader) (core.Packet, error) {
	httpResp, err := http.ReadResponse(reader, nil)
	if err != nil {
		return nil, err
	}

	httpResp.Header.Set("Connection", "keep-alive")
	httpResp.Header.Set("Keep-Alive", "timeout=4")
	httpResp.Close = false
	return &httpResponsePacket{httpResp}, nil
}
