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
	reader := bufio.NewReader(srm.svr)
	req, err := srm.decoder.decode(reader)
	if err != nil {
		countlog.Error("event!server.failed to read request", "err", err)
		return
	}
	qualifier := core.Route(req)
	clt, err := client.Connect(qualifier)
	if err != nil {
		countlog.Error("event!server.failed to connect client", "err", err)
		return
	}
	err = req.Write(clt)
	if err != nil {
		clt.Close()
		countlog.Error("event!server.failed to write request", "err", err)
		return
	}
	// clt is handed over to this goroutine now
	srm.forwardResponsesInGoroutine(clt)
	srm.forwardFollowingRequests(reader, clt)
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
		defer clt.Close() // ensure clt is closed after io.Copy finished
		written, err := io.Copy(srm.svr, clt)
		countlog.Debug("event!server.copied response", "written", written, "err", err)
	}()
}

func (srm *stream) forwardFollowingRequests(reader *bufio.Reader, clt client.OutboundClient) {
	for {
		req, err := srm.decoder.decode(reader)
		if err != nil {
			countlog.Error("event!server.failed to read request", "err", err)
			return
		}
		err = req.Write(clt)
		if err != nil {
			countlog.Error("event!server.failed to write request", "err", err)
			return
		}
	}
}
