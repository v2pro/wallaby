package main

import (
	"net/http"
	"runtime"
)

var oneKB []byte
var tenKB []byte

func init() {
	for i := 0; i < 1024; i++ {
		oneKB = append(oneKB, 'A')
	}
	for i := 0; i < 1024 * 10; i++ {
		tenKB = append(tenKB, 'A')
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	http.HandleFunc("/oneKB", func(w http.ResponseWriter, req *http.Request) {
		w.Write(oneKB)
	})
	http.HandleFunc("/tenKB", func(w http.ResponseWriter, req *http.Request) {
		w.Write(tenKB)
	})
	http.ListenAndServe("127.0.0.1:8849", nil)
}