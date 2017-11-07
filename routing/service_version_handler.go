package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/v2pro/plz/countlog"
	"io"
	"net/http"
	"strconv"
	"time"
)

func list(versions *ServiceVersions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var vs []ServiceVersion
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

func set(versions *ServiceVersions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decode := json.NewDecoder(r.Body)
		var s ServiceVersion
		err := decode.Decode(&s)
		if err != nil {
			countlog.Error("event!Version", "set", err)
			w.WriteHeader(501)
			return
		}

		if ok := versions.Set(&s); !ok {
			w.WriteHeader(501)
		} else {
			w.WriteHeader(200)
		}
	}
}

func get(versions *ServiceVersions) func(w http.ResponseWriter, r *http.Request) {
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

func del(versions *ServiceVersions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decode := json.NewDecoder(r.Body)
		var s ServiceVersion
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

func helloHandler(buildTimeStamp int) func(w http.ResponseWriter, r *http.Request) {
	var timestampStr = strconv.Itoa(buildTimeStamp)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Build-Timestamp", timestampStr)
	}
}

type InboundService struct {
	addr           string
	mux            *http.ServeMux
	server         *http.Server
	version        *ServiceVersions
	buildTimeStamp int
}

func NewInboundService(addr string,
	thisVersions *ServiceVersions,
	buildTimeStamp int) *InboundService {

	mux := http.NewServeMux()

	// return proxy build timestamp
	mux.HandleFunc("/hello", helloHandler(buildTimeStamp))

	mux.HandleFunc("/list", list(thisVersions))
	mux.HandleFunc("/set", set(thisVersions))
	mux.HandleFunc("/get", get(thisVersions))
	mux.HandleFunc("/del", del(thisVersions))

	inboundService := &InboundService{
		addr:           addr,
		mux:            mux,
		server:         &http.Server{Addr: addr, Handler: mux},
		version:        thisVersions,
		buildTimeStamp: buildTimeStamp,
	}

	return inboundService
}

func (s *InboundService) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			return
		}
	}()
}

func (s *InboundService) Shutdown() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return s.server.Shutdown(ctx)
}
