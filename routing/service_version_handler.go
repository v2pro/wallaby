package version

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/routing"
	"io"
	"net/http"
	"strconv"
	"time"
)

func list(versions *routing.ServiceVersions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var vs []routing.ServiceVersion
		for _, v := range versions.List() {
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
}

func set(versions *routing.ServiceVersions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decode := json.NewDecoder(r.Body)
		var s routing.ServiceVersion
		err := decode.Decode(&s)
		if err != nil {
			fmt.Printf("error %v", err)
			w.WriteHeader(501)
			return
		}
		fmt.Printf("input %v", s)

		if ok := versions.Set(&s); !ok {
			w.WriteHeader(501)
		} else {
			w.WriteHeader(200)
		}
	}
}
func get(versions *routing.ServiceVersions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s := versions.Get()
		if s == nil {
			w.WriteHeader(204)
		} else if data, err := json.Marshal(*s); err != nil {
			w.WriteHeader(501)
		} else {
			w.WriteHeader(200)
			io.WriteString(w, string(data))
		}
	}
}
func del(versions *routing.ServiceVersions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decode := json.NewDecoder(r.Body)
		var s routing.ServiceVersion
		err := decode.Decode(&s)
		if err != nil {
			fmt.Printf("error %v", err)
			w.WriteHeader(501)
			return
		}
		fmt.Printf("input %v", s)

		if ok := versions.Del(s.Address); !ok {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(200)
		}
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, r.URL.Path)
}

type InboundService struct {
	port    int
	mux     *http.ServeMux
	server  *http.Server
	version *routing.ServiceVersions
}

func NewInboundService(port int, versionFilePath string) *InboundService {
	thisVersions := routing.NewServiceVersions(versionFilePath)
	if thisVersions.Start() != nil {
		panic("start thisVersions fail")
		countlog.Error("event!NewInboundService", "New ", "start thisVersions fail")
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/echo", echoHandler)

	mux.HandleFunc("/list", list(thisVersions))
	mux.HandleFunc("/set", set(thisVersions))
	mux.HandleFunc("/get", get(thisVersions))
	mux.HandleFunc("/del", del(thisVersions))

	return &InboundService{
		port:    port,
		mux:     mux,
		server:  &http.Server{Addr: ":" + strconv.Itoa(port), Handler: mux},
		version: thisVersions,
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
