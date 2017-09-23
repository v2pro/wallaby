package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveService(t *testing.T) {
	s, err := ResolveService("coupon-localhost@default")
	assert.Nil(t, err)
	fmt.Printf("ResolveService = %s\n", s)
}
