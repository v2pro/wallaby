package supervisor

import (
	"context"
	"encoding/json"
	"github.com/v2pro/plz/countlog"
	"io"
	"net/http"
	"strconv"
	"time"
)

func list(procMgr *ProcMgr) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var proc_list []ProcInfo = procMgr.List()
		data, err := json.Marshal(proc_list)
		if err != nil {
			w.WriteHeader(501)
		} else {
			w.WriteHeader(200)
			io.WriteString(w, string(data))
		}
	}
}

func start(procMgr *ProcMgr) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decode := json.NewDecoder(r.Body)
		var p ProcInfo
		err := decode.Decode(&p)
		if err != nil {
			countlog.Error("event!Supervisord", "start", err)
			w.WriteHeader(501)
			return
		}
		err = procMgr.StartProc(p)
		if err != nil {
			countlog.Error("event!Supervisord", "start", err)
			w.WriteHeader(501)
		} else {
			w.WriteHeader(200)
		}
	}
}

func stop(procMgr *ProcMgr) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	}
}

func clean(procMgr *ProcMgr) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	}
}

func helloHandler(buildTimeStamp int) func(w http.ResponseWriter, r *http.Request) {
	var timestampStr = strconv.Itoa(buildTimeStamp)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Build-Timestamp", timestampStr)
	}
}

type Supervisord struct {
	addr           string
	mux            *http.ServeMux
	server         *http.Server
	procMgr        *ProcMgr
	buildTimeStamp int
}

func NewSupervisord(addr string,
	procMgr *ProcMgr,
	buildTimeStamp int) *Supervisord {

	mux := http.NewServeMux()

	// return proxy build timestamp
	mux.HandleFunc("/hello", helloHandler(buildTimeStamp))

	mux.HandleFunc("/list", list(procMgr))
	mux.HandleFunc("/start", start(procMgr))
	mux.HandleFunc("/stop", stop(procMgr))
	mux.HandleFunc("/clean", clean(procMgr))

	supervisord := &Supervisord{
		addr:           addr,
		mux:            mux,
		server:         &http.Server{Addr: addr, Handler: mux},
		procMgr:        procMgr,
		buildTimeStamp: buildTimeStamp,
	}

	return supervisord
}

func (s *Supervisord) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			return
		}
	}()
}

func (s *Supervisord) Shutdown() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.procMgr.StopAll()
	return s.server.Shutdown(ctx)
}
