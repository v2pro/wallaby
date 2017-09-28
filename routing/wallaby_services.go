package routing

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/core/coretype"
	"github.com/v2pro/wallaby/countlog"
)

var (
	defaultRouteTable = map[string]string{}
	nodeList          = []*ServiceNode{}
)

type WallabyServiceDiscover struct {
}

func (w WallabyServiceDiscover) FindServiceKindAddr(sk *core.ServiceKind) (*net.TCPAddr, error) {
	addr, err := FindServiceKindAddr(sk)
	if err != nil {
		return nil, err
	}
	return net.ResolveTCPAddr("tcp", addr)
}

func GetCurrentVersion() string {
	return nodeList[0].Version
}

func GetCurrentServiceNode() *ServiceNode {
	return nodeList[0]
}

func GetNextVersion() string {
	if len(nodeList) > 0 {
		return nodeList[1].Version
	} else {
		return ""
	}
}

func FindServiceKindAddr(sk *core.ServiceKind) (string, error) {
	//addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8005")
	if addr, ok := defaultRouteTable[sk.String()]; ok {
		return addr, nil
	}
	return "", fmt.Errorf("FindServiceAddr fail to find %s", sk)
}

func init() {
	if !InitRouteTable() {
		panic("InitRouteTable fail")
	}
}

func InitRouteTable() bool {
	namingFile, err := os.Open(GetRoot() + "/wallaby-services.json")
	if err != nil {
		countlog.Errorf("no wallaby-services.json found: %s", err.Error())
		return false
	}

	jsonParser := json.NewDecoder(namingFile)
	if err = jsonParser.Decode(&nodeList); err != nil {
		countlog.Errorf("fail to parse config file: %s", err.Error())
		return false
	}
	if len(nodeList) == 0 {
		countlog.Error("no services found")
		return false
	}

	for _, node := range nodeList {
		serviceNode, err := ResolveService(node)
		if err != nil {
			countlog.Error(err.Error())
			continue
		}
		defaultRouteTable[serviceNode.String()] = node.Address
	}
	if len(nodeList) == 0 {
		countlog.Error("no services found")
		return false
	}
	return true
}

func GetRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	relRoot := path.Dir(filename) + "/.."
	root, err := filepath.Abs(relRoot)
	if err != nil {
		panic(fmt.Sprintf("Wrong path: %s @ %s\n", err.Error(), relRoot))
	}
	return root
}

type ServiceNode struct {
	Service  string `json:"service"`
	Cluster  string `json:"cluster"`
	Version  string `json:"version"`
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
}

func ResolveService(node *ServiceNode) (*core.ServiceKind, error) {
	s := &core.ServiceKind{}
	s.Name = node.Service
	s.Cluster = node.Cluster
	s.Version = node.Version
	s.Protocol = coretype.Protocol(node.Protocol)
	return s, nil
}
