package main

import (
	"testing"
	"net"
	"net/http"
	"bytes"
	"net/http/httputil"
	"bufio"
	"io/ioutil"
)

func Benchmark_long_connection(b *testing.B) {
	addr := "127.0.0.1:8849"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		b.Error(err)
		return
	}
	req, err := http.NewRequest("GET", "/oneKB", bytes.NewBufferString(""))
	if err != nil {
		b.Error(err)
		return
	}
	req.Host = addr
	req.Header.Set("Connection", "keep-alive")
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		b.Error(err)
		return
	}
	reader := bufio.NewReaderSize(conn, 2048)
	for i := 0; i < b.N; i++ {
		_, err := conn.Write(reqBytes)
		if err != nil {
			b.Error(err)
			return
		}
		resp, err := http.ReadResponse(reader, nil)
		if err != nil {
			b.Error(err)
			return
		}
		if resp.Body != nil {
			_, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				b.Error(err)
				return
			}
		}
	}
}
