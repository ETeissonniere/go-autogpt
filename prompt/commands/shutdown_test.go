package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShutdownCommand(t *testing.T) {
	ret, err := (&ShutdownCommand{}).Execute(nil)
	assert.Equal(t, ErrShutdown, err)
	assert.Equal(t, "", ret)
}
