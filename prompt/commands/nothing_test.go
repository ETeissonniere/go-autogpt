package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoNothingCommand(t *testing.T) {
	ret, err := (&DoNothingCommand{}).Execute(nil)
	assert.Nil(t, err)
	assert.Equal(t, "nothing done", ret)
}
