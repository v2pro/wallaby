package main

import (
	"github.com/v2pro/wallaby/server"
	"runtime"
	"github.com/v2pro/plz/countlog"
)

func main() {
	runtime.GOMAXPROCS(1)
	asyncLogWriter := countlog.NewAsyncLogWriter(
		countlog.LEVEL_INFO,
		countlog.NewFileLogOutput("STDERR"))
	asyncLogWriter.Start()
	countlog.LogWriters = append(countlog.LogWriters, asyncLogWriter)
	server.Start()
}