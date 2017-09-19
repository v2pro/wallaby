package server

import (
	"net/http"
	"bufio"
	"github.com/v2pro/wallaby/core"
)

type httpInboundRequest struct {
	*http.Request
}

func (req *httpInboundRequest) Feature() map[string]string {
	return map[string]string{}
}

type httpRequestDecoder struct {
}

func (decoder *httpRequestDecoder) decode(reader *bufio.Reader) (core.InboundRequest, error){
	httpReq, err := http.ReadRequest(reader)
	if err != nil {
		return nil, err
	}
	return &httpInboundRequest{httpReq}, nil
}