package version

import (
	rand2 "math/rand"
)

const (
	Started  = string("Started")
	Running  = string("Running")
	Stopped  = string("Stopped")
	Failed   = string("Failed")
	ChanSize = 10
)

type (
	ServiceVersion struct {
		Address  string `json:"address"` //`comment host:port`
		Name     string `json:"name"`    //`comment service name`
		Version  string `json:"version"`
		Status   string `json:"status"`
		Priority uint32 `json:"priority"` //`comment more priority more requests, zero stands for no request`
	}

	SetHandler func(old_sv, new_sv *ServiceVersion)
	DelHandler func(sv *ServiceVersion)

	ServiceVersions struct {
		versionList   []*ServiceVersion `comment key is Address`
		totalPriority uint32
		listReq       chan listRequest `comment list request channel`
		setReq        chan setRequest  `comment set request channel`
		getReq        chan getRequest  `comment get request channel`
		delReq        chan delRequest  `comment del request channel`
		stop          chan struct{}    `comment stop action channel`
		setHandler    SetHandler
		delHandler    DelHandler
	}

	listRequest struct {
		res chan []*ServiceVersion
	}

	setRequest struct {
		sv  ServiceVersion `comment copy ServerVersion`
		res chan bool
	}

	getRequest struct {
		res chan *ServiceVersion
	}

	delRequest struct {
		address string
		res     chan bool
	}
)

func (s *ServiceVersions) SetHandler(h SetHandler) {
	s.setHandler = h
}

func (s *ServiceVersions) DelHandler(h DelHandler) {
	s.delHandler = h
}

func (s *ServiceVersions) list() []*ServiceVersion {
	return s.versionList
}

func (s *ServiceVersions) set(sv *ServiceVersion) {
	var old_sv *ServiceVersion = nil
	var idx int = -1
	for k, v := range s.versionList {
		if v.Address == sv.Address {
			old_sv = s.versionList[k]
			idx = k
			break
		}
	}
	if s.setHandler != nil {
		s.setHandler(old_sv, sv)
	}
	if old_sv != nil {
		// replace
		if old_sv.Status == Running {
			s.totalPriority -= old_sv.Priority
		}
		s.versionList[idx] = sv
	} else {
		// add
		s.versionList = append(s.versionList, sv)
	}
	if sv.Status == Running {
		s.totalPriority += sv.Priority
	}
}

func (s *ServiceVersions) get() *ServiceVersion {
	if s.totalPriority == 0 {
		return nil
	}
	var r uint32 = uint32(rand2.Int31()) % s.totalPriority
	for _, v := range s.versionList {
		if v.Status == Running {
			if r < v.Priority {
				return v
			} else {
				r -= v.Priority
			}
		}
	}
	return nil
}

func (s *ServiceVersions) del(address string) bool {
	var sv *ServiceVersion
	for k, v := range s.versionList {
		if v.Address == address {
			if v.Status == Running {
				s.totalPriority -= v.Priority
			}
			sv = s.versionList[k]
			s.versionList[k] = s.versionList[len(s.versionList)-1]
			s.versionList = s.versionList[:len(s.versionList)-1]
			if s.delHandler != nil {
				s.delHandler(sv)
			}
			return true
		}
	}
	return false
}

func (s *ServiceVersions) List() []*ServiceVersion {
	req := listRequest{res: make(chan []*ServiceVersion)}
	s.listReq <- req
	return <-req.res
}

func (s *ServiceVersions) Set(sv *ServiceVersion) bool {
	if sv == nil {
		return false
	}
	req := setRequest{sv: *sv, res: make(chan bool)}
	s.setReq <- req
	return <-req.res
}

func (s *ServiceVersions) Get() *ServiceVersion {
	req := getRequest{res: make(chan *ServiceVersion)}
	s.getReq <- req
	return <-req.res
}

func (s *ServiceVersions) Del(Add string) bool {
	req := delRequest{address: Add, res: make(chan bool)}
	s.delReq <- req
	return <-req.res
}

func (s *ServiceVersions) Start() {
	go func() {
		for {
			select {
			case <-s.stop:
				return
			case req := <-s.listReq:
				req.res <- s.list()
			case req := <-s.getReq:
				req.res <- s.get()
			case req := <-s.setReq:
				s.set(&req.sv)
				req.res <- true
			case req := <-s.delReq:
				req.res <- s.del(req.address)
			}
		}
	}()
}

func (s *ServiceVersions) Stop() {
	s.stop <- struct{}{}
}

func NewServiceVersions() *ServiceVersions {
	var sv *ServiceVersions = &ServiceVersions{
		versionList:   []*ServiceVersion{},
		totalPriority: 0,
		listReq:       make(chan listRequest, ChanSize),
		setReq:        make(chan setRequest, ChanSize),
		getReq:        make(chan getRequest, ChanSize),
		delReq:        make(chan delRequest, ChanSize),
		stop:          make(chan struct{}),
		setHandler:    nil,
		delHandler:    nil,
	}
	return sv
}

var serviceVersions *ServiceVersions

func init() {
	serviceVersions = NewServiceVersions()
}

func GetServiceVersions() *ServiceVersions {
	return serviceVersions
}
