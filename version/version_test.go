package version

import (
	"github.com/v2pro/wallaby/util"
	"strconv"
	"sync"
	"testing"
)

func TestServiceVersions(t *testing.T) {
	svs := NewServiceVersions()
	util.AssertEqual(t, len(svs.list()), 0, "empty")
	svs.Start()
	util.AssertEqual(t, len(svs.List()), 0, "empty")
	util.AssertEqual(t, svs.totalPriority, uint32(0), "0 totalPriority")

	sv1 := &ServiceVersion{
		Address:  "127.0.0.1:1",
		Status:   Running,
		Priority: 100,
	}
	util.AssertEqual(t, svs.Set(sv1), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 1, "empty")
	util.AssertEqual(t, svs.totalPriority, uint32(100), "100 totalPriority")

	sv2 := &ServiceVersion{
		Address:  "127.0.0.1:1",
		Status:   Running,
		Priority: 200,
	}
	util.AssertEqual(t, svs.Set(sv2), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 1, "empty")
	util.AssertEqual(t, svs.totalPriority, uint32(200), "200 totalPriority")

	sv3 := &ServiceVersion{
		Address:  "127.0.0.1:2",
		Status:   Running,
		Priority: 300,
	}

	util.AssertEqual(t, svs.Set(sv3), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 2, "empty")
	util.AssertEqual(t, svs.totalPriority, uint32(500), "500 totalPriority")

	sv4 := &ServiceVersion{
		Address:  "127.0.0.1:2",
		Status:   Stopped,
		Priority: 300,
	}

	util.AssertEqual(t, svs.Set(sv4), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 2, "empty")
	util.AssertEqual(t, svs.totalPriority, uint32(200), "200 totalPriority")
	svs.Stop()
}

func TestGetServiceVersions(t *testing.T) {
	svs := NewServiceVersions()
	svs.Start()
	util.AssertEqual(t, len(svs.List()), 0, "empty")
	sv := &ServiceVersion{
		Address:  "127.0.0.1:1",
		Status:   Running,
		Priority: 100,
	}

	util.AssertEqual(t, svs.Set(sv), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 1, "empty")
	util.AssertEqual(t, svs.totalPriority, uint32(100), "100 totalPriority")
	util.AssertEqual(t, svs.Get().Priority, uint32(100), "100 totalPriority")
	sv.Priority = 0
	util.AssertEqual(t, svs.Get().Priority, uint32(100), "100 totalPriority")

	sv = &ServiceVersion{
		Address:  "127.0.0.1:2",
		Status:   Running,
		Priority: 100,
	}
	util.AssertEqual(t, svs.Set(sv), true, "set failed")
	util.AssertEqual(t, svs.totalPriority, uint32(200), "200 totalPriority")
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
		Priority: 100,
	}
	util.AssertEqual(t, svs.Set(sv), true, "set failed")
	util.AssertEqual(t, svs.totalPriority, uint32(100), "100 totalPriority")
	count = map[string]int{}
	for i := 0; i < 1000; i++ {
		sv := svs.Get()
		count[sv.Address] += 1
	}
	util.AssertEqual(t, len(count), 1, "len count")
	t.Log(count)

	util.AssertEqual(t, svs.Del(sv.Address), true, "set failed")
	util.AssertEqual(t, len(svs.List()), 1, "1")
	util.AssertEqual(t, svs.totalPriority, uint32(100), "100 totalPriority")
	util.AssertEqual(t, svs.Del(sv.Address), false, "set failed")

	svs.Stop()

}
func perf(s *ServiceVersions, wg *sync.WaitGroup, idx int) {
	wg.Add(1)
	defer wg.Done()
	var loop int = 1e3
	for i := 0; i < loop; i = i + 1 {
		sv := &ServiceVersion{
			Address:  "127.0.0.1:" + strconv.Itoa(i),
			Status:   Stopped,
			Priority: 100,
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
	var wg sync.WaitGroup

	svs := NewServiceVersions()
	svs.Start()
	for i := 0; i < 10; i = i + 1 {
		go perf(svs, &wg, i)
	}
	wg.Wait()
	util.AssertEqual(t, svs.totalPriority, uint32(0), "0 totalPriority")
	util.AssertEqual(t, len(svs.List()), 0, "1")

	svs.Stop()
}

func TestCallback(t *testing.T) {
	svs := NewServiceVersions()
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
		Priority: 200,
	}
	util.AssertEqual(t, svs.Set(sv), true, "set failed")
	util.AssertEqual(t, s_count, 1, "set count 1")
	util.AssertEqual(t, svs.Del(sv.Address), true, "set failed")
	util.AssertEqual(t, d_count, 1, "set count 1")
}
