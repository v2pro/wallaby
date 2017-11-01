package routing

import (
	"encoding/json"
	"github.com/v2pro/wallaby/core/codec"
	"github.com/v2pro/wallaby/datacenter"
	"github.com/v2pro/wallaby/util"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"testing"
)

var empty_packet codec.Packet

func TestServiceVersionsLoadStore(t *testing.T) {

	svs := NewServiceVersions("tmp/file/not/exists")
	if svs.Start() == nil {
		panic("start fail")
	}

	store_file := "/tmp/0"
	svs = NewServiceVersions(store_file)
	defer os.Remove(store_file)
	if svs.Start() != nil {
		panic("start fail")
	}
	util.AssertEqual(t, len(svs.list()), 0, "empty")

	sv1 := &ServiceVersion{
		Address:  "127.0.0.1:1",
		Status:   Running,
		Value:    100,
		Operator: datacenter.OperatorRandom,
	}
	t.Log("empty svs.list() ", svs.list())
	util.AssertEqual(t, svs.Set(sv1), true, "set failed")
	util.AssertEqual(t, len(svs.list()), 1, "not empty")
	t.Log("one svs.list() ", svs.list())
	cat, err := exec.Command("cat", store_file).Output()
	t.Log("cat file ", string(cat[:]), err)
	svs.Stop()

	svs2 := NewServiceVersions(store_file)
	if svs2.Start() != nil {
		panic("start fail")
	}
	util.AssertEqual(t, svs2.Set(sv1), true, "set failed")
	util.AssertEqual(t, len(svs2.list()), 1, "not empty")
	t.Log("one svs.list() ", svs2.list())
	svs2.Stop()
}

func TestServiceVersions(t *testing.T) {
	svs := NewServiceVersions("/tmp/1")
	defer os.Remove("/tmp/1")
	util.AssertEqual(t, len(svs.list()), 0, "empty")
	if svs.Start() != nil {
		panic("start fail")
	}
	util.AssertEqual(t, len(svs.List()), 0, "empty")
	util.AssertEqual(t, svs.totalPriority, int32(0), "0 totalPriority")

	sv1 := &ServiceVersion{
		Address:  "127.0.0.1:1",
		Status:   Running,
		Value:    100,
		Operator: datacenter.OperatorRandom,
	}
	util.AssertEqual(t, svs.Set(sv1), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 1, "empty")
	util.AssertEqual(t, svs.totalPriority, int32(100), "100 totalPriority")

	sv2 := &ServiceVersion{
		Address:  "127.0.0.1:1",
		Status:   Running,
		Value:    200,
		Operator: datacenter.OperatorRandom,
	}
	util.AssertEqual(t, svs.Set(sv2), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 1, "empty")
	util.AssertEqual(t, svs.totalPriority, int32(200), "200 totalPriority")

	sv3 := &ServiceVersion{
		Address:  "127.0.0.1:2",
		Status:   Running,
		Value:    300,
		Operator: datacenter.OperatorRandom,
	}

	util.AssertEqual(t, svs.Set(sv3), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 2, "empty")
	util.AssertEqual(t, svs.totalPriority, int32(500), "500 totalPriority")

	sv4 := &ServiceVersion{
		Address:  "127.0.0.1:2",
		Status:   Stopped,
		Value:    300,
		Operator: datacenter.OperatorRandom,
	}
	bin, err := json.Marshal(sv4)
	t.Log("json.Marshal", string(bin[:]), err)

	util.AssertEqual(t, svs.Set(sv4), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 2, "empty")
	util.AssertEqual(t, svs.totalPriority, int32(200), "200 totalPriority")
	svs.Stop()
}

func TestGetServiceVersions(t *testing.T) {
	svs := NewServiceVersions("/tmp/2")
	defer os.Remove("/tmp/2")
	svs.Start()
	util.AssertEqual(t, len(svs.List()), 0, "empty")
	sv := &ServiceVersion{
		Address:  "127.0.0.1:1",
		Status:   Running,
		Value:    100,
		Operator: datacenter.OperatorRandom,
	}

	util.AssertEqual(t, svs.Set(sv), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 1, "empty")
	util.AssertEqual(t, svs.totalPriority, int32(100), "100 totalPriority")
	util.AssertEqual(t, svs.Route(empty_packet).Value, int32(100), "100 totalPriority")
	sv.Value = 0
	util.AssertEqual(t, svs.Route(empty_packet).Value, int32(100), "100 totalPriority")

	sv = &ServiceVersion{
		Address:  "127.0.0.1:2",
		Status:   Running,
		Value:    100,
		Operator: datacenter.OperatorRandom,
	}
	util.AssertEqual(t, svs.Set(sv), true, "set failed")
	util.AssertEqual(t, svs.totalPriority, int32(200), "200 totalPriority")
	var count map[string]int = map[string]int{}
	for i := 0; i < 1000; i++ {
		sv := svs.Get()
		count[sv.Address] += 1
	}
	util.AssertEqual(t, len(count), 2, "len count")
	t.Log(count)

	sv = &ServiceVersion{
		Address:  "127.0.0.1:2",
		Status:   Stopped,
		Value:    100,
		Operator: datacenter.OperatorRandom,
	}
	util.AssertEqual(t, svs.Set(sv), true, "set failed")
	util.AssertEqual(t, svs.totalPriority, int32(100), "100 totalPriority")
	count = map[string]int{}
	for i := 0; i < 1000; i++ {
		sv := svs.Get()
		count[sv.Address] += 1
	}
	util.AssertEqual(t, len(count), 1, "len count")
	t.Log(count)

	util.AssertEqual(t, svs.Del(sv.Address), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 1, "1")
	util.AssertEqual(t, svs.totalPriority, int32(100), "100 totalPriority")
	util.AssertEqual(t, svs.Del(sv.Address), false, "set failed")

	svs.Stop()

}

func perf(s *ServiceVersions, wg *sync.WaitGroup, loop int) {
	defer wg.Done()
	for i := 0; i < loop; i = i + 1 {
		sv := &ServiceVersion{
			Address:  "127.0.0.1:" + strconv.Itoa(i),
			Status:   Stopped,
			Value:    100,
			Operator: datacenter.OperatorRandom,
		}
		if i < 3 {
			s.Set(sv)
		}
		for j := 0; j < 10; j = j + 1 {
			s.Get()
		}
		if i < 3 {
			s.Del(sv.Address)
		}
	}
}

func TestPerf(t *testing.T) {
	// 100w per second
	var loop int = 1 //e3
	var wg sync.WaitGroup

	svs := NewServiceVersions("/tmp/3")
	defer os.Remove("/tmp/3")
	svs.Start()
	wg.Add(10)
	for i := 0; i < 10; i = i + 1 {
		go perf(svs, &wg, loop)
	}
	wg.Wait()
	util.AssertEqual(t, svs.totalPriority, int32(0), "0 totalPriority")
	util.AssertEqual(t, len(svs.List()), 0, "1")

	svs.Stop()
}

func TestCallback(t *testing.T) {
	svs := NewServiceVersions("/tmp/4")
	defer os.Remove("/tmp/4")
	svs.Start()
	defer svs.Stop()

	s_count := 0
	sh := func(old, new *ServiceVersion) {
		s_count = s_count + 1
	}

	d_count := 0
	dh := func(old *ServiceVersion) {
		d_count = d_count + 1
	}
	svs.SetHandler(sh)
	svs.DelHandler(dh)

	sv := &ServiceVersion{
		Address:  "127.0.0.1:1",
		Status:   Running,
		Value:    200,
		Operator: datacenter.OperatorRandom,
	}
	util.AssertEqual(t, svs.Set(sv), true, "set failed")
	util.AssertEqual(t, s_count, 1, "set count 1")
	util.AssertEqual(t, svs.Del(sv.Address), true, "set failed")
	util.AssertEqual(t, d_count, 1, "set count 1")
}
