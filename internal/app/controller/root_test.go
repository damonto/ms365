package controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrap(t *testing.T) {
	code, wrapped := rootCtl.wrap(200, "message")

	assert.Equal(t, code, 200)
	assert.Equal(t, wrapped.Message, "message")

	_, wrapped = rootCtl.wrap(200, map[string]string{
		"testing": "testing",
	})

	assert.Equal(t, wrapped.Data.(map[string]string)["testing"], "testing")
}
