package server

import (
	"bytes"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/config"
	"github.com/v2pro/wallaby/server"
	"net/http"
	"os"
	"runtime"
	"testing"
	"time"
)

var oneKB []byte
var tenKB []byte

var sv1_addr string
var sv2_addr string

var set_sv1_json []byte
var set_sv2_json []byte

func init() {
	runtime.GOMAXPROCS(1)

	for i := 0; i < 1024; i++ {
		oneKB = append(oneKB, 'A')
	}
	for i := 0; i < 1024*10; i++ {
		tenKB = append(tenKB, 'A')
	}

	sv1_addr = "127.0.0.1:8851"
	sv2_addr = "127.0.0.1:8852"

	set_sv1_json = []byte(`{"address" : "` + sv1_addr + `", "name" : "test", "version" : "1.0.3", "status" : "Running", "value" : 10, "operator" : "random"}`)
	set_sv2_json = []byte(`{"address" : "` + sv2_addr + `", "name" : "test", "version" : "1.0.2", "status" : "Running", "value" : 10, "operator" : "random"}`)
}

func StartEchoServer(addr string) *http.Server {
	countlog.Info("event!StartEchoServer", "Addr", addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/echo", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		// allow pre-flight headers
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
		writer.Header().Set("Origin-Addr", addr)
		request.Write(writer)
	})
	server := &http.Server{Addr: addr, Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil {
		}
	}()
	return server
}

func StartProxyServer() *server.ProxyServer {
	asyncLogWriter := countlog.NewAsyncLogWriter(
		countlog.LEVEL_INFO,
		//countlog.LEVEL_DEBUG,
		countlog.NewFileLogOutput("STDERR"))
	asyncLogWriter.Start()
	countlog.LogWriters = append(countlog.LogWriters, asyncLogWriter)
	proxy := server.ProxyServer{}
	go func() {
		if err := proxy.Start(); err != nil {
			countlog.Error("event!server.failed to accept outbound", "err", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	return &proxy
}

func SetProxyVersion(sv_json []byte) {
	client := &http.Client{}
	host := "http://" + config.VersionHandlerAddr

	req, _ := http.NewRequest("GET", host+"/set", bytes.NewBuffer(sv_json))
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic("error StatusCode")
	}
}

func Echo(host string, r []byte) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", host+"/echo", bytes.NewBuffer(r))
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		panic("set sv1 fail")
	}
	return resp
}

func TestProxyEchoTest(t *testing.T) {
	proxy := StartProxyServer()
	defer os.Remove(config.ProxyServiceVersionConfig)
	defer proxy.Stop()

	sv1 := StartEchoServer(sv1_addr)
	defer sv1.Shutdown(nil)

	sv2 := StartEchoServer(sv2_addr)
	defer sv2.Shutdown(nil)

	SetProxyVersion(set_sv1_json)
	SetProxyVersion(set_sv2_json)
	/*
		b.ResetTimer()
	*/

	host := "http://" + config.ProxyAddr

	for i := 0; i < 1e3; i++ {
		Echo(host, oneKB)
		//t.Logf("resp header %v", resp.Header["Origin-Addr"])
	}

	return

}
