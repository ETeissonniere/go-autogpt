package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShutdownCommand(t *testing.T) {
	ret, err := (&ShutdownCommand{}).Execute(nil)
	assert.Nil(t, err)
	assert.Equal(t, "", ret)
}
