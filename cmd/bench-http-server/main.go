package main

import (
	"net/http"
	"runtime"
)

var oneKB []byte

func init() {
	for i := 0; i < 1024; i++ {
		oneKB = append(oneKB, 'A')
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	http.HandleFunc("/oneKB", func(w http.ResponseWriter, req *http.Request) {
		w.Write(oneKB)
	})
	http.ListenAndServe("127.0.0.1:8849", nil)
}