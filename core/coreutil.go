package core

import (
	"fmt"
	"strings"
)

// ResolveService qualifier is like "coupon-localhost@default"
func ResolveService(qualifier string) (*ServiceKind, error) {
	s := &ServiceKind{}
	parts := strings.Split(qualifier, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("ResolveService fail: count not match||qualifier=%s", qualifier)
	}
	s.Name = parts[0]
	parts = strings.Split(parts[1], "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("ResolveService fail: count not match||qualifier=%s", qualifier)
	}
	s.Cluster = parts[0]
	s.Version = parts[1]
	s.Protocol = Http
	return s, nil
}
