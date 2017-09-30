package infra

import (
	"encoding/json"
	"github.com/v2pro/wallaby/countlog"
)

func JsonEncode(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(str)
}

// JsonDecode if input is json format，return a map
// if input is empty, wrong format, return nil
// can't handle json array，e.g., [1,2,3]，use JsonDecodeArr instead
func JsonDecode(encodeStr string) map[string]interface{} {
	var j interface{}
	if parseErr := json.Unmarshal([]byte(encodeStr), &j); parseErr != nil {
		//dlog.Warningf("JsonDecode err: %s, json string: %s", parseErr, encodeStr)
		return nil
	}
	switch ret := j.(type) {
	case map[string]interface{}:
		return ret
	case string, int, int64, float64:
		return nil
	default:
		countlog.Errorf("JsonDecode(%s) return %T(%v), expect map[string]interface{}", encodeStr, ret, ret)
		return nil
	}
}

// handle json array，e.g., [1,2,3]
func JsonDecodeArr(encodeStr string) []interface{} {
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
		countlog.Errorf("JsonDecode(%s) return %T(%v), expect []interface{}", encodeStr, ret, ret)
		return nil
	}
}