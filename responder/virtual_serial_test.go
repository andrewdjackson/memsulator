package responder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVirtualSerialPort_Create(t *testing.T) {
	var err error

	vs := NewVirtualSerialPort()
	err = vs.CreateVirtualPorts()

	assert.Equal(t, err, nil)
}
