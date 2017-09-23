package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	ProxyPort = 8080
)

// this is just a demo ReverseProxy using the built-in golang lib.
func main() {
	//go runServerDemo(8866)
	mapUrlPrefixToVirtualHost("/this", "http://127.0.0.1:8006")
	mapHostnameToVirtualHost("this.service:8080/", "http://127.0.0.1:8006")

	//go runServerDemo(8867)
	mapUrlPrefixToVirtualHost("/new", "http://127.0.0.1:8005")
	mapHostnameToVirtualHost("new.service:8080/", "http://127.0.0.1:8005")

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "proxy status: OK")
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", ProxyPort), nil)
	if err != nil {
		panic(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}

func runServerDemo(port int) {
	serverInstMux := http.NewServeMux()
	serverInstMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world! -- port: %d", port)
	})
	serverInst := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: serverInstMux}
	serverInst.ListenAndServe()
}

func mapUrlPrefixToVirtualHost(prefix string, serverAddr string) {
	vhost1, err := url.Parse(serverAddr)
	if err != nil {
		panic(err)
	}
	proxy1 := httputil.NewSingleHostReverseProxy(vhost1)
	http.HandleFunc(prefix, handler(proxy1))
}

func mapHostnameToVirtualHost(prefix string, serverAddr string) {
	vhost1, err := url.Parse(serverAddr)
	if err != nil {
		panic(err)
	}
	proxy1 := httputil.NewSingleHostReverseProxy(vhost1)
	http.HandleFunc(prefix, handler(proxy1))
}
