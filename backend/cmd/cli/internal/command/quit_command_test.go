package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuitCommand(t *testing.T) {
	quit := NewQuitCommand()
	q, err := quit.Run([]string{})
	assert.NoError(t, err)
	assert.True(t, q)
}
