package server

import (
	"net"
	"github.com/v2pro/wallaby/countlog"
	"bufio"
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/client"
	"io"
	"os"
)

type decoder interface {
	decodeRequest(reader *bufio.Reader) (core.Packet, error)
	decodeResponse(reader *bufio.Reader) (core.Packet, error)
}

type stream struct {
	svr     *net.TCPConn
	decoder decoder
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
	requestReader := bufio.NewReaderSize(srm.svr, 2048)
	responseReader := bufio.NewReaderSize(os.Stdin, 2048) // reader not set yet
	for srm.roundtrip(requestReader, responseReader) {

	}
}

func (srm *stream) roundtrip(requestReader *bufio.Reader, responseReader *bufio.Reader) bool {
	req, err := srm.decoder.decodeRequest(requestReader)
	if err == io.EOF {
		countlog.Trace("event!server.inbound conn closed")
		return false
	}
	if err != nil {
		countlog.Warn("event!server.failed to read request", "err", err)
		return false
	}
	qualifier := core.Route(req)
	clt, err := client.Get(qualifier)
	if err != nil {
		countlog.Warn("event!server.failed to connect client", "err", err)
		return false
	}
	if srm.handleRequest(req, clt, responseReader) {
		return true
	}
	clt, err = client.GetNew(qualifier)
	if err != nil {
		countlog.Warn("event!server.failed to re-connect client", "err", err)
		return false
	}
	return srm.handleRequest(req, clt, responseReader)
}

func (srm *stream) handleRequest(req core.Packet, clt client.OutboundClient, responseReader *bufio.Reader) bool {
	defer clt.Close()
	responseReader.Reset(clt)
	err := req.Write(clt)
	if err != nil {
		countlog.Warn("event!server.failed to write request", "err", err)
		return false
	}
	resp, err := srm.decoder.decodeResponse(responseReader)
	if err != nil {
		countlog.Warn("event!server.failed to read response", "err", err)
		return false
	}
	err = resp.Write(srm.svr)
	if err != nil {
		countlog.Warn("event!server.failed to write response", "err", err)
		return false
	}
	return true
}
