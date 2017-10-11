package infra

import (
	"encoding/json"
	"fmt"
	"github.com/v2pro/plz/countlog"
)

// JSONEncode return encoded json string if no error, return "" if error occurs
func JSONEncode(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(str)
}

// JSONDecode if input is json format，return a map
// if input is empty, wrong format, return nil
// can't handle json array，e.g., [1,2,3]，use JSONDecodeArr instead
func JSONDecode(encodeStr string) map[string]interface{} {
	var j interface{}
	if parseErr := json.Unmarshal([]byte(encodeStr), &j); parseErr != nil {
		//dlog.Warningf("JSONDecode err: %s, json string: %s", parseErr, encodeStr)
		return nil
	}
	switch ret := j.(type) {
	case map[string]interface{}:
		return ret
	case string, int, int64, float64:
		return nil
	default:
		countlog.Error(fmt.Sprintf("JSONDecode(%s) return %T(%v), expect map[string]interface{}", encodeStr, ret, ret))
		return nil
	}
}

// JSONDecodeArr handle json array，e.g., [1,2,3]
func JSONDecodeArr(encodeStr string) []interface{} {
	var j interface{}
	if parseErr := json.Unmarshal([]byte(encodeStr), &j); parseErr != nil {
		return nil
	}
	switch ret := j.(type) {
	case []interface{}:
		return ret
	case string, int, int64, float64:
		return nil
	default:
		countlog.Error(fmt.Sprintf("JSONDecode(%s) return %T(%v), expect []interface{}", encodeStr, ret, ret))
		return nil
	}
}
