package main

import (
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/server"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)
	asyncLogWriter := countlog.NewAsyncLogWriter(
		countlog.LEVEL_INFO,
		countlog.NewFileLogOutput("STDERR"))
	asyncLogWriter.LogFormatter = &countlog.CompactFormat{}
	asyncLogWriter.Start()
	countlog.LogWriters = append(countlog.LogWriters, asyncLogWriter)
	proxy := server.ProxyServer{}
	proxy.Start()
}
