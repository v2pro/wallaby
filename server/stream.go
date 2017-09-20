package server

import (
	"net"
	"github.com/v2pro/wallaby/countlog"
	"bufio"
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/client"
	"io"
)

type requestDecoder interface {
	decode(reader *bufio.Reader) (core.InboundRequest, error)
}

type stream struct {
	svr     *net.TCPConn
	decoder requestDecoder
}

func (srm *stream) proxy() {
	defer func() {
		recovered := recover()
		if recovered != nil {
			countlog.Fatal("event!server.panic", "err", recovered,
				"stacktrace", countlog.ProvideStacktrace)
		}
	}()
	defer srm.svr.Close()
	c := &capture{
		reader: srm.svr,
		bytesRead: make([]byte, 0, 2048),
	}
	reader := bufio.NewReaderSize(c, 2048)
	req, err := srm.decoder.decode(reader)
	if err != nil {
		countlog.Warn("event!server.failed to read request", "err", err)
		return
	}
	qualifier := core.Route(req)
	clt, err := client.Connect(qualifier)
	if err != nil {
		countlog.Warn("event!server.failed to connect client", "err", err)
		return
	}
	defer clt.Close()
	_, err = clt.Write(c.bytesRead)
	if err != nil {
		countlog.Debug("event!server.failed to write request", "err", err)
		return
	}
	srm.forwardResponsesInGoroutine(clt)
	srm.forwardFollowingRequests(clt)
}

func (srm *stream) forwardResponsesInGoroutine(clt client.OutboundClient) {
	go func() {
		defer func() {
			recovered := recover()
			if recovered != nil {
				countlog.Fatal("event!server.panic", "err", recovered,
					"stacktrace", countlog.ProvideStacktrace)
			}
		}()
		defer srm.svr.Close()
		defer clt.Close()
		written, err := io.CopyBuffer(srm.svr, clt, make([]byte, 2048))
		countlog.Debug("event!server.copied response", "written", written, "err", err)
	}()
}

func (srm *stream) forwardFollowingRequests(clt client.OutboundClient) {
	written, err := io.CopyBuffer(clt, srm.svr, make([]byte, 2048))
	countlog.Debug("event!server.copied request", "written", written, "err", err)
}
