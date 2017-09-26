package server

import (
	"github.com/v2pro/wallaby/client"
	"github.com/v2pro/wallaby/config"
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/core/codec"
	"github.com/v2pro/wallaby/countlog"
	"io"
	"net"
	"time"
)

type stream struct {
	svr        *net.TCPConn
	svrCodec   codec.Codec
	svrCapture *codec.Capture
	cltCapture *codec.Capture
}

func newStream(svr *net.TCPConn, svrCodec codec.Codec) *stream {
	svrCapture := &codec.Capture{}
	svrCapture.SetReader(svr)
	cltCapture := &codec.Capture{}
	return &stream{
		svr:        svr,
		svrCodec:   svrCodec,
		svrCapture: svrCapture,
		cltCapture: cltCapture,
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
	target := core.HowToRoute(&core.ServerRequest{
		Packet: req,
	})
	clt, err := client.Get(target.ServiceInstance)
	if err != nil {
		countlog.Warn("event!server.failed to connect client", "err", err)
		return false
	}
	if srm.handleRequest(clt, req) {
		return true
	}
	// because the client from pool might be disconnected
	// we get a "new" client and try again
	// the "old" client will be discarded because read/write incurred error which marked it as invalid
	// this way we can expire invalid client and re-fill the pool with new one
	countlog.Debug("event!server.re-connect client")
	clt, err = client.GetNew(target.ServiceInstance.ServiceKind)
	if err != nil {
		countlog.Warn("event!server.failed to re-connect client", "err", err)
		return false
	}
	return srm.handleRequest(clt, req)
}

func (srm *stream) readRequest() codec.Packet {
	for {
		srm.svr.SetReadDeadline(time.Now().Add(time.Duration(int(time.Second) * config.ClientReadTimeout)))
		req, err := srm.svrCodec.DecodeRequest(srm.svrCapture)
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

func (srm *stream) handleRequest(clt client.Client, req codec.Packet) bool {
	defer clt.Close()
	resp, err := clt.Handle(req, srm.cltCapture)
	if err != nil {
		return false
	}
	err = srm.svrCodec.EncodeResponse(resp, srm.svr)
	if err != nil {
		countlog.Warn("event!server.failed to write response", "err", err)
		return false
	}
	return true
}
