package datacenter

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/util/infra"
	"regexp"
	"strconv"
)

const (
	// OperatorEq equal
	OperatorEq = "="

	// OperatorGt greater than
	OperatorGt = ">"

	// OperatorGe greater equal than
	OperatorGe = ">="

	// OperatorLt less than
	OperatorLt = "<"

	// OperatorLe less equal than
	OperatorLe = "<="

	// OperatorRegex use regular expression to match
	OperatorRegex = "regex"

	// OperatorRandom random value from 0 - 99
	OperatorRandom = "random"
)

var (
	routingSetting = &RoutingSetting{}
	validOperators = map[string]bool{}
)

func init() {
	validOperators[OperatorEq] = true
	validOperators[OperatorGt] = true
	validOperators[OperatorGe] = true
	validOperators[OperatorLt] = true
	validOperators[OperatorLe] = true
	validOperators[OperatorRegex] = true
	validOperators[OperatorRandom] = true

	ReadUserSetting()
}

// GetRoutingSetting get config from user-settings.json
func GetRoutingSetting() *RoutingSetting {
	if routingSetting.IsEmpty() {
		ReadUserSetting()
	}
	return routingSetting
}

// ReadUserSetting loads config from user-settings.json
func ReadUserSetting() bool {
	settingFile, err := os.Open(getRoot() + "/user-settings.json")
	if err != nil {
		countlog.Error("event!no user-settings.json found", "err", err)
		return false
	}

	jsonParser := json.NewDecoder(settingFile)
	if err = jsonParser.Decode(routingSetting); err != nil {
		countlog.Error("event!fail to parse user-settings file", "err", err)
		return false
	}
	if !routingSetting.IsValid() {
		countlog.Error("event!user-settings is invalid", "routingSetting", infra.JSONEncode(routingSetting))
		return false
	}
	return true
}

// getRoot get root path of project directory
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

// RoutingSetting is json struct for user-settings.json rule items
type RoutingSetting struct {
	Hashkey  string `json:"hashkey"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// NewRoutingSetting return a new RoutingSetting object
func NewRoutingSetting(hashkey string, operator string, val string) *RoutingSetting {
	return &RoutingSetting{
		Hashkey: hashkey, Operator: operator, Value: val,
	}
}

// IsEmpty checks if any value is empty
func (rs RoutingSetting) IsEmpty() bool {
	if rs.Hashkey == "" {
		return true
	}
	if rs.Operator == "" {
		return true
	}
	if rs.Value == "" {
		return true
	}
	return false
}

// IsValid return true if values are not empty and operator is valid, return false otherwise
func (rs RoutingSetting) IsValid() bool {
	if rs.IsEmpty() {
		return false
	}
	if !isOperatorValid(rs.Operator) {
		return false
	}
	return true
}

// RunRoutingRule return true if the value satisfies the rule, return false otherwise
func (rs RoutingSetting) RunRoutingRule(hashVal string) bool {
	switch rs.Operator {
	case OperatorEq:
		return hashVal == rs.Value
	case OperatorGe:
		return hashVal >= rs.Value
	case OperatorGt:
		return hashVal > rs.Value
	case OperatorLt:
		return hashVal < rs.Value
	case OperatorLe:
		return hashVal <= rs.Value
	case OperatorRegex:
		var validValue = regexp.MustCompile(rs.Value)
		return validValue.MatchString(hashVal)
	case OperatorRandom:
		hashInt, err := strconv.Atoi(rs.Value)
		if err != nil {
			countlog.Error("event!fail to parse RoutingSetting.Value to int", "err", err)
			return false
		}
		return infra.RandomPercent() >= hashInt
	default:
		return false
	}
}

func isOperatorValid(operator string) bool {
	if _, ok := validOperators[operator]; ok {
		return true
	}
	return false
}
