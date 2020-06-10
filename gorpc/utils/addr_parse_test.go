package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseServicePath(t *testing.T) {
	_,_,err :=ParseServicePath("Hello,World")
	assert.NotNil(t,err)

	_,_,err  =ParseServicePath("Hello/World")
	assert.NotNil(t,err)
	serviceName, method, err := ParseServicePath("/Hello/World")
	assert.Equal(t, serviceName, "Hello")
	assert.Equal(t, method, "World")
	assert.Equal(t, err, nil)
}