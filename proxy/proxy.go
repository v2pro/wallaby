package proxy

import (
	"net"
	"github.com/v2pro/wallaby/countlog"
	"fmt"
)

func Start() {
	addr := "127.0.0.1:8848"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		countlog.Error("event!proxy.failed to bind proxy port", "err", err)
		return
	}
	countlog.Info("event!proxy.started", "addr", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			countlog.Error("event!proxy.failed to accept outbound", "err", err)
			return
		}
		go handleInbound(conn.(*net.TCPConn))
	}
}

func handleInbound(conn *net.TCPConn) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			countlog.Fatal("event!proxy.panic", "err", recovered,
				"stacktrace", countlog.ProvideStacktrace)
		}
	}()
	fmt.Println("!!!")
}