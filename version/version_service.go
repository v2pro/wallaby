package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var (
	thisVersions *ServiceVersions
)

func init() {
	thisVersions = GetServiceVersions()
	thisVersions.Start()
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, r.URL.Path)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	var vs []ServiceVersion
	for _, v := range thisVersions.List() {
		vs = append(vs, *v)
	}
	data, err := json.Marshal(vs)
	if err != nil {
		w.WriteHeader(501)
	} else {
		w.WriteHeader(200)
		io.WriteString(w, string(data))
	}
}

func setHandler(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	var s ServiceVersion
	err := decode.Decode(&s)
	if err != nil {
		fmt.Printf("error %v", err)
		w.WriteHeader(501)
		return
	}
	fmt.Printf("input %v", s)

	if ok := thisVersions.Set(&s); !ok {
		w.WriteHeader(501)
	} else {
		w.WriteHeader(200)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	s := thisVersions.Get()
	if s == nil {
		w.WriteHeader(204)
	} else if data, err := json.Marshal(*s); err != nil {
		w.WriteHeader(501)
	} else {
		w.WriteHeader(200)
		io.WriteString(w, string(data))
	}
}

func delHandler(w http.ResponseWriter, r *http.Request) {
	decode := json.NewDecoder(r.Body)
	var s ServiceVersion
	err := decode.Decode(&s)
	if err != nil {
		fmt.Printf("error %v", err)
		w.WriteHeader(501)
		return
	}
	fmt.Printf("input %v", s)

	if ok := thisVersions.Del(s.Address); !ok {
		w.WriteHeader(204)
	} else {
		w.WriteHeader(200)
	}
}

func registerHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/echo", echoHandler)

	mux.HandleFunc("/list", listHandler)
	mux.HandleFunc("/set", setHandler)
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/del", delHandler)

	return mux
}

type InboundService struct {
	port   int
	mux    *http.ServeMux
	server *http.Server
}

func NewInboundService(port int) *InboundService {
	mux := registerHandler()
	return &InboundService{
		port:   port,
		mux:    mux,
		server: &http.Server{Addr: ":" + strconv.Itoa(port), Handler: mux},
	}
}

func (s *InboundService) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			return
		}
	}()
}

func (s *InboundService) Shutdown() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.server.Shutdown(ctx)
}
