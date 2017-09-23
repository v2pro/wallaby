package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/v2pro/wallaby/countlog"
)

func FindServiceAddr(qualifier *ServiceKind) (string, error) {
	//addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8005")
	if addr, ok := defaultRouteTable[qualifier.String()]; ok {
		return addr, nil
	}
	return "", fmt.Errorf("FindServiceAddr fail to find %s", qualifier)
}

func HowToRoute(serverConn *ServerConn) RoutingMode {
	return ""
}

func Route(serverRequest *ServerRequest) *RoutingDecision {
	return &RoutingDecision{
		ServiceInstance: &ServiceInstance{
			ServiceKind: serviceNodeList[0],
		},
	}
}

var (
	defaultRouteTable = map[string]string{}
	serviceNodeList   = []*ServiceKind{}
)

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

func init() {
	if !InitRouteTable() {
		panic("InitRouteTable fail")
	}
}

func InitRouteTable() bool {
	namingFile, err := os.Open(GetRoot() + "/naming.json")
	if err != nil {
		countlog.Errorf("no naming.json found: %s", err.Error())
		return false
	}
	nodeList := []*ServiceNode{}
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
		serviceNode, err := ResolveService(node.Qualifier)
		if err != nil {
			countlog.Error(err.Error())
			continue
		}
		serviceNodeList = append(serviceNodeList, serviceNode)
		defaultRouteTable[node.Qualifier] = node.Address
	}
	if len(nodeList) == 0 {
		countlog.Error("no services found")
		return false
	}
	return true
}
