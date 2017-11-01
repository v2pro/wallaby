package routing

import (
	"encoding/json"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/core/codec"
	"github.com/v2pro/wallaby/core/coretype"
	"github.com/v2pro/wallaby/datacenter"
	"io/ioutil"
	rand2 "math/rand"
	"os"
	"sort"
	"strconv"
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
		Address  string            `json:"address"` //`comment host:port`
		Service  string            `json:"name"`    //`comment service name`
		Version  string            `json:"version"`
		PWD      string            `json:"pwd"`
		Protocol coretype.Protocol `json:"protocol"`
		Tag      string            `json:"tag"`
		Status   string            `json:"status"`
		HashKey  string            `json:"hashkey"`
		Operator string            `json:"operator"`
		Value    int32             `json:"value"` //`comment more priority more requests, zero stands for no request`
		Rule     datacenter.RoutingSetting
	}

	VersionList []*ServiceVersion

	SetHandler func(old_sv, new_sv *ServiceVersion)
	DelHandler func(sv *ServiceVersion)

	ServiceVersions struct {
		filepath      string            `comment serialize file`
		versionList   VersionList       `comment key is Address`
		totalPriority int32             `comment valid for random operator`
		listReq       chan listRequest  `comment list request channel`
		setReq        chan setRequest   `comment set request channel`
		getReq        chan getRequest   `comment get request channel`
		routeReq      chan routeRequest `comment route request channel`
		delReq        chan delRequest   `comment del request channel`
		stop          chan struct{}     `comment stop action channel`
		setHandler    SetHandler
		delHandler    DelHandler
		serviceName   string
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

	routeRequest struct {
		packet codec.Packet
		res    chan *ServiceVersion
	}

	delRequest struct {
		address string
		res     chan bool
	}
)

func (l VersionList) Len() int {
	return len(l)
}

func (l VersionList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l VersionList) Less(i, j int) bool {
	// desc
	return l[i].Version > l[j].Version
}

func (s *ServiceVersion) GenRule() {
	s.Rule = datacenter.RoutingSetting{
		Hashkey: s.HashKey, Operator: s.Operator, Value: strconv.Itoa(int(s.Value)),
	}
}

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
	defer func() {
		if err := s.store(); err != nil {
			countlog.Error("event!ServiceVersions.set", "store ", err)
			panic(err)
		}
	}()
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
		if old_sv.Status == Running && old_sv.Rule.Operator == datacenter.OperatorRandom {
			s.totalPriority -= old_sv.Value
		}
		s.versionList[idx] = sv
	} else {
		// add
		s.versionList = append(s.versionList, sv)
	}
	if sv.Status == Running && sv.Rule.Operator == datacenter.OperatorRandom {
		s.totalPriority += sv.Value
	}

	// sort versionList by timestamp desc
	sort.Sort(s.versionList)
}

func (s *ServiceVersions) del(address string) bool {
	defer func() {
		if err := s.store(); err != nil {
			countlog.Error("event!ServiceVersions.set", "store ", err)
			panic(err)
		}
	}()

	var sv *ServiceVersion
	for k, v := range s.versionList {
		if v.Address == address {
			if v.Status == Running && v.Rule.Operator == datacenter.OperatorRandom {
				s.totalPriority -= v.Value
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

func (s *ServiceVersions) get() *ServiceVersion {
	if s.totalPriority == 0 {
		return nil
	}
	var r int32 = rand2.Int31() % s.totalPriority
	for _, v := range s.versionList {
		if v.Status == Running {
			if r < v.Value {
				return v
			} else {
				r -= v.Value
			}
		}
	}
	return nil
}

// version is sorted desc, compute route rule from latest one to oldest one
// if no rule matched, use the oldest version
func (s *ServiceVersions) route(packet codec.Packet) *ServiceVersion {
	var r int32 = 0
	if s.totalPriority != 0 {
		r = rand2.Int31() % s.totalPriority
	}
	var default_version *ServiceVersion = nil
	for _, v := range s.versionList {
		if v.Status != Running {
			continue
		}
		default_version = v
		if v.Rule.IsValid() {
			if v.Rule.Operator != datacenter.OperatorRandom {
				if packet != nil && v.Rule.RunRoutingRule((packet).GetFeature(v.Rule.Hashkey)) {
					return v
				}
			} else {
				if r < v.Value {
					return v
				} else {
					r -= v.Value
				}
			}
		}
	}
	return default_version
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
	sv.GenRule()
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

func (s *ServiceVersions) Route(packet codec.Packet) *ServiceVersion {
	req := routeRequest{packet: packet, res: make(chan *ServiceVersion)}
	s.routeReq <- req
	return <-req.res
}

func (s *ServiceVersions) load() error {
	// TODO: if no file exists, empty config
	f, ferr := os.OpenFile(s.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if ferr == nil {
		f.Close()
	} else {
		countlog.Warn("event!ServiceVersions.load", s.filepath, ferr)
		return ferr
	}

	bin, err := ioutil.ReadFile(s.filepath)
	if err != nil {
		countlog.Error("event!ServiceVersions.load", s.filepath, err)
		return err
	}
	countlog.Info("event!ServiceVersions.load", s.filepath, "{"+string(bin[:])+"}")

	if len(bin) == 0 {
		countlog.Warn("event!ServiceVersions.load", s.filepath, "empty file")
		return nil
	}
	err = json.Unmarshal(bin, &s.versionList)
	if err != nil {
		countlog.Error("event!ServiceVersions.store", "write file ", bin, err)
		return err
	}
	return nil
}

func (s *ServiceVersions) store() error {
	bin, err := json.Marshal(s.versionList)
	countlog.Debug("event!ServiceVersions.store", "json", bin)
	if err != nil {
		countlog.Error("event!ServiceVersions.store", "json marshal ", bin, err)
		return err
	}
	err = ioutil.WriteFile(s.filepath, bin, 0644)
	if err != nil {
		countlog.Error("event!ServiceVersions.store", "write file ", bin, err)
		return err
	}
	return nil
}

func (s *ServiceVersions) Start() error {
	// TODO reentrance
	// check load / store
	err := s.load()
	if err != nil {
		countlog.Error("event!ServiceVersions.Start", "load ", err)
		return err
	}
	err = s.store()
	if err != nil {
		countlog.Error("event!ServiceVersions.Start", "store ", err)
		return err
	}
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
			case req := <-s.routeReq:
				req.res <- s.route(req.packet)
			}
		}
	}()
	return nil
}

func (s *ServiceVersions) Stop() {
	s.stop <- struct{}{}
}

func NewServiceVersions(filePath string) *ServiceVersions {
	var sv *ServiceVersions = &ServiceVersions{
		filepath:      filePath,
		versionList:   []*ServiceVersion{},
		totalPriority: 0,
		listReq:       make(chan listRequest, ChanSize),
		setReq:        make(chan setRequest, ChanSize),
		getReq:        make(chan getRequest, ChanSize),
		delReq:        make(chan delRequest, ChanSize),
		routeReq:      make(chan routeRequest, ChanSize),
		stop:          make(chan struct{}),
		setHandler:    nil,
		delHandler:    nil,
	}
	return sv
}
