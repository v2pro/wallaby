package server

import (
	"net"
	"github.com/v2pro/wallaby/countlog"
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/client"
	"io"
	"time"
	"github.com/v2pro/wallaby/core/codec"
)

type stream struct {
	svr             *net.TCPConn
	serverDecoder   codec.Decoder
	requestCapture  *codec.Capture
	responseCapture *codec.Capture
}

func newStream(svr *net.TCPConn, serverDecoder codec.Decoder) *stream {
	requestCapture := &codec.Capture{}
	requestCapture.Reset(svr)
	return &stream{
		svr:             svr,
		serverDecoder:   serverDecoder,
		requestCapture:  requestCapture,
		responseCapture: &codec.Capture{},
	}
}

func (srm *stream) proxy() {
	defer func() {
		recovered := recover()
		if recovered != nil {
			countlog.Fatal("event!server.panic", "err", recovered,
				"stacktrace", countlog.ProvideStacktrace)
		}
		srm.svr.Close()
	}()
	for srm.roundtrip() {

	}
}

func (srm *stream) roundtrip() bool {
	req := srm.readRequest()
	if req == nil {
		return false
	}
	qualifier, clientDecoder := core.Route(req)
	clt, err := client.Get(qualifier)
	if err != nil {
		countlog.Warn("event!server.failed to connect client", "err", err)
		return false
	}
	if srm.handleRequest(clt, req, clientDecoder) {
		return true
	}
	countlog.Debug("event!server.re-connect client")
	clt, err = client.GetNew(qualifier)
	if err != nil {
		countlog.Warn("event!server.failed to re-connect client", "err", err)
		return false
	}
	return srm.handleRequest(clt, req, clientDecoder)
}

func (srm *stream) readRequest() codec.Packet {
	for {
		srm.svr.SetReadDeadline(time.Now().Add(time.Second * 5))
		req, err := srm.serverDecoder.DecodeRequest(srm.requestCapture)
		if err == io.EOF {
			countlog.Trace("event!server.inbound conn closed")
			return nil
		}
		if err != nil {
			countlog.Warn("event!server.failed to read request", "err", err)
			return nil
		}
		return req
	}
}

func (srm *stream) handleRequest(
	clt client.Client, req codec.Packet, clientDecoder codec.Decoder) bool {
	defer clt.Close()
	resp, err := clt.Handle(req, srm.responseCapture, clientDecoder)
	if err != nil {
		return false
	}
	err = resp.Write(srm.svr)
	if err != nil {
		countlog.Warn("event!server.failed to write response", "err", err)
		return false
	}
	return true
}
