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
	OperatorEq     = "="
	OperatorGt     = ">"
	OperatorGe     = ">="
	OperatorLt     = "<"
	OperatorLe     = "<="
	OperatorRegex  = "regex"
	OperatorRandom = "random" // 0 - 99
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
}

func GetRoutingSetting() *RoutingSetting {
	if routingSetting.IsEmpty() {
		ReadUserSetting()
	}
	return routingSetting
}

func isOperatorValid(operator string) bool {
	if _, ok := validOperators[operator]; ok {
		return true
	}
	return false
}

func init() {
	ReadUserSetting()
}

func ReadUserSetting() bool {
	settingFile, err := os.Open(GetRoot() + "/user-settings.json")
	if err != nil {
		countlog.Errorf("no user-settings.json found: %s", err.Error())
		return false
	}

	jsonParser := json.NewDecoder(settingFile)
	if err = jsonParser.Decode(routingSetting); err != nil {
		countlog.Errorf("fail to parse user-settings file: %s", err.Error())
		return false
	}
	if !routingSetting.IsValid() {
		countlog.Errorf("user-settings is invalid: %s", infra.JsonEncode(routingSetting))
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

type RoutingSetting struct {
	Hashkey  string `json:"hashkey"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

func NewRoutingSetting(hashkey string, operator string, val string) *RoutingSetting {
	return &RoutingSetting{
		Hashkey: hashkey, Operator: operator, Value: val,
	}
}
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

func (rs RoutingSetting) IsValid() bool {
	if rs.IsEmpty() {
		return false
	}
	if !isOperatorValid(rs.Operator) {
		return false
	}
	return true
}

func (rs RoutingSetting) RunRoutingRule(hashVal string) bool {
	switch rs.Operator {
	case OperatorEq:
		if hashVal == rs.Value {
			return true
		}
		return false
	case OperatorGe:
		if hashVal >= rs.Value {
			return true
		}
		return false
	case OperatorGt:
		if hashVal > rs.Value {
			return true
		}
		return false
	case OperatorLt:
		if hashVal < rs.Value {
			return true
		}
		return false
	case OperatorLe:
		if hashVal <= rs.Value {
			return true
		}
		return false
	case OperatorRegex:
		var validValue = regexp.MustCompile(rs.Value)
		return validValue.MatchString(hashVal)
	case OperatorRandom:
		hashInt, err := strconv.Atoi(rs.Value)
		if err != nil {
			countlog.Errorf("fail to parse RoutingSetting.Value to int: %s", err)
			return false
		}
		return infra.RandomPercent() >= hashInt
	default:
		return false
	}
}
