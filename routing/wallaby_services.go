package routing

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/core"
	"github.com/v2pro/wallaby/core/coretype"
)

var (
	defaultRouteTable = map[string]string{}
	nodeList          = []*ServiceNode{}
)

// GetCurrentVersion get the version string of current running service
func GetCurrentVersion() string {
	return nodeList[0].Version
}

// GetCurrentServiceNode get the current running service info
func GetCurrentServiceNode() *ServiceNode {
	return nodeList[0]
}

// GetNextVersion get the new version string of service
func GetNextVersion() string {
	if len(nodeList) > 0 {
		return nodeList[1].Version
	}
	return ""
}

// FindServiceKindAddr return the ip+port string for the given ServiceKind
func FindServiceKindAddr(sk *core.ServiceKind) (string, error) {
	//addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8005")
	if addr, ok := defaultRouteTable[sk.Qualifier()]; ok {
		return addr, nil
	}
	return "", fmt.Errorf("FindServiceAddr fail to find %s", sk)
}

func init() {
	if !GetWallabyServices() {
		panic("GetWallabyServices fail")
	}
}

// GetWallabyServices load values from wallaby-services.json
func GetWallabyServices() bool {
	namingFile, err := os.Open(getRoot() + "/wallaby-services.json")
	if err != nil {
		countlog.Error("event!no wallaby-services.json found", "err", err)
		return false
	}

	jsonParser := json.NewDecoder(namingFile)
	if err = jsonParser.Decode(&nodeList); err != nil {
		countlog.Error("event!fail to parse config file", "err", err)
		return false
	}
	if len(nodeList) == 0 {
		countlog.Error("event!no services found")
		return false
	}

	for _, node := range nodeList {
		serviceNode, err := resolveService(node)
		if err != nil {
			countlog.Error(err.Error())
			continue
		}
		defaultRouteTable[serviceNode.Qualifier()] = node.Address
	}
	if len(defaultRouteTable) == 0 {
		countlog.Error("no services found")
		return false
	}
	return true
}

func getRoot() string {
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

// ServiceNode is the json struct for wallaby-services.json items
type ServiceNode struct {
	Service  string `json:"service"`
	Cluster  string `json:"cluster"`
	Version  string `json:"version"`
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
	Tag      string `json:"tag"`
}

func resolveService(node *ServiceNode) (*core.ServiceKind, error) {
	s := &core.ServiceKind{}
	s.Name = node.Service
	s.Cluster = node.Cluster
	s.Version = node.Version
	s.Protocol = coretype.Protocol(node.Protocol)
	s.Tag = node.Tag
	return s, nil
}
