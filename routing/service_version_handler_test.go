package routing

import (
	"bytes"
	"encoding/json"
	"github.com/v2pro/wallaby/util"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	addr := "127.0.0.1:18022"
	filePath := "test.json"
	defer os.Remove(filePath)
	thisVersions := NewServiceVersions(filePath)
	if thisVersions.Start() != nil {
		panic("start thisVersions fail")
	}
	server := NewInboundService(addr, thisVersions, 1)
	server.Start()
	defer server.Shutdown()
	time.Sleep(100 * time.Millisecond)

	client := &http.Client{}
	host := "http://" + addr
	// empty list
	req, _ := http.NewRequest("GET", host+"/list", nil)
	resp, err := client.Do(req)
	//defer req.Body.Close()
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Logf("list %v", resp)
		t.Logf("body %v", string(body))

		util.AssertEqual(t, 200, resp.StatusCode, "get 200")
		util.AssertEqual(t, string(body), "null", "null body")
		decode := json.NewDecoder(resp.Body)
		var s []*ServiceVersion
		err := decode.Decode(s)
		util.AssertNotEqual(t, err, nil, "null body")
		util.AssertEqual(t, len(s), 0, "null json")
	} else {
		t.Logf("%v", err)
		util.AssertEqual(t, 1, 0, "fail")
	}

	// get nil
	req, _ = http.NewRequest("GET", host+"/get", nil)
	resp, err = client.Do(req)
	if err == nil {
		t.Logf("get %v", resp)
		util.AssertEqual(t, 204, resp.StatusCode, "get 200")
		var vs ServiceVersion
		decode := json.NewDecoder(resp.Body)
		err := decode.Decode(&vs)
		util.AssertNotEqual(t, err, nil, "not nil")
	} else {
		t.Logf("%v", err)
		util.AssertEqual(t, 1, 0, "fail")
	}

	// set
	set_json := []byte(`{"address" : "1", "name" : "test", "version" : "1.0.0", "pwd":"/home/work/1_0_0", "status" : "Running", "value" : 10, "operator" : "random"}`)
	t.Logf("%v", string(set_json))
	req, _ = http.NewRequest("GET", host+"/set", bytes.NewBuffer(set_json))
	resp, err = client.Do(req)
	defer req.Body.Close()
	if err == nil {
		util.AssertEqual(t, 200, resp.StatusCode, "get 200")
	} else {
		panic("set fail")
	}

	// list again
	req, _ = http.NewRequest("GET", host+"/list", nil)
	resp, err = client.Do(req)
	if err == nil {
		//body, _ := ioutil.ReadAll(resp.Body)
		t.Logf("list %v", resp)
		//t.Logf("list body %v", string(body))
		util.AssertEqual(t, 200, resp.StatusCode, "get 200")
		//util.AssertNotEqual(t, string(body), "null", "null body")
		var vs []ServiceVersion
		//err := json.Unmarshal(body, vs)
		decode := json.NewDecoder(resp.Body)
		err := decode.Decode(&vs)
		if err != nil {
			t.Logf("json decode error %f", err)
		}
		util.AssertEqual(t, len(vs), 1, "one addr")
		util.AssertEqual(t, vs[0].Status, Running, "running status")
		util.AssertEqual(t, vs[0].PWD, "/home/work/1_0_0", "running status")

	} else {
		t.Logf("%v", err)
		util.AssertEqual(t, 1, 0, "fail")
	}

	// get one
	req, _ = http.NewRequest("GET", host+"/get", nil)
	resp, err = client.Do(req)
	if err == nil {
		t.Logf("get %v", resp)
		util.AssertEqual(t, 200, resp.StatusCode, "get 200")
		var vs ServiceVersion
		decode := json.NewDecoder(resp.Body)
		err := decode.Decode(&vs)
		if err != nil {
			t.Logf("json decode error %f", err)
		}
		util.AssertEqual(t, vs.Address, "1", "null json")
	} else {
		t.Logf("%v", err)
		util.AssertEqual(t, 1, 0, "fail")
	}

	// del
	set_json = []byte(`{"address" : "1", "name" : "test", "version" : "1.0.0", "status" : "Running", "priority" : 10}`)
	req, _ = http.NewRequest("GET", host+"/del", bytes.NewBuffer(set_json))
	resp, err = client.Do(req)
	if err == nil {
		t.Logf("get %v", resp)
		util.AssertEqual(t, 200, resp.StatusCode, "get 200")
	} else {
		t.Logf("%v", err)
		util.AssertEqual(t, 1, 0, "fail")
	}

	// get nil
	req, _ = http.NewRequest("GET", host+"/get", nil)
	resp, err = client.Do(req)
	if err == nil {
		t.Logf("get %v", resp)
		util.AssertEqual(t, 204, resp.StatusCode, "get 200")
		var vs ServiceVersion
		decode := json.NewDecoder(resp.Body)
		err := decode.Decode(&vs)
		util.AssertNotEqual(t, err, nil, "not nil")
	} else {
		t.Logf("%v", err)
		util.AssertEqual(t, 1, 0, "fail")
	}
}
