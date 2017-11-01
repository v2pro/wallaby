package routing

import (
	"github.com/v2pro/wallaby/datacenter"
	"os"
	"testing"
)

type TestPacket struct {
	feature map[string]string
}

func (p *TestPacket) GetFeature(key string) string {
	value, ok := p.feature[key]
	if ok {
		return value
	} else {
		return ""
	}
}

var config_file = "version_strategy_test.json"

func CreateRandomServiceVersionStrategy() *VersionRoutingStrategy {

	var vrs = NewVersionRoutingStrategy("test", config_file, "127.0.0.1:12345")
	if vrs == nil {
		panic("NewVersionRoutingStrategy fail")
	}

	var svs *ServiceVersions = vrs.ServiceVersions()

	sv1 := &ServiceVersion{
		Address:  "127.0.0.1:1",
		Status:   Running,
		Value:    100,
		Operator: datacenter.OperatorRandom,
	}
	svs.Set(sv1)

	sv2 := &ServiceVersion{
		Address:  "127.0.0.1:2",
		Status:   Running,
		Value:    100,
		Operator: datacenter.OperatorRandom,
	}
	svs.Set(sv2)

	if len(svs.list()) != 2 {
		panic("set service version fail")
	}
	return vrs
}

func TestRandomRouting(t *testing.T) {
	vrs := CreateRandomServiceVersionStrategy()
	defer os.Remove(config_file)
	packet := TestPacket{}
	var addrCounter = map[string]int{}
	for i := 0; i < 1000; i++ {
		var ver *ServiceVersion = vrs.Route(&packet)
		//var ver *ServiceVersion = vrs.ServiceVersions().Get()

		addrCounter[ver.Address] += 1
	}
	t.Log(addrCounter)
	return
}

func BenchmarkRandomRouting(b *testing.B) {
	vrs := CreateRandomServiceVersionStrategy()
	defer os.Remove(config_file)
	packet := TestPacket{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vrs.Route(&packet)
	}
	return
}

func TestRegexRouting(t *testing.T) {
	return
}
