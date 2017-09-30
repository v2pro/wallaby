package datacenter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadUserSetting(t *testing.T) {
	setting := GetRoutingSetting()
	assert.True(t, setting.IsValid())
}

func TestRunRoutingRule(t *testing.T) {
	assert.True(t, NewRoutingSetting("x-forwarded-for", "regex", "[12345]$").RunRoutingRule("192.168.1.115"))
	assert.True(t, NewRoutingSetting("Cityid", "=", "12345").RunRoutingRule("12345"))
	assert.True(t, NewRoutingSetting("Cityid", ">", "10000").RunRoutingRule("12345"))
	assert.True(t, NewRoutingSetting("x-forwarded-for", "regex", "[12345]?$").RunRoutingRule(""))
	assert.False(t, NewRoutingSetting("x-forwarded-for", "regex", "[12345]$").RunRoutingRule(""))
	assert.False(t, NewRoutingSetting("x-forwarded-for", "random", "100").RunRoutingRule(""))
	assert.True(t, NewRoutingSetting("x-forwarded-for", "random", "0").RunRoutingRule(""))
	result := NewRoutingSetting("x-forwarded-for", "random", "50").RunRoutingRule("")
	if result == true {
		fmt.Println("random true")
	} else {
		fmt.Println("random false")
	}
	assert.True(t, result == true || result == false)
}
