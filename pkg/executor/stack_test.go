package executor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack(t *testing.T) {
	s := newStack[int]()

	s.push(1)
	assert.False(t, s.empty())

	s.push(2)
	n := s.pop()
	assert.Equal(t, n, 2)
	assert.False(t, s.empty())

	n = s.pop()
	assert.Equal(t, n, 1)
	assert.True(t, s.empty())

	s.push(3)
	s.push(4)
	s.push(5)

	assert.Equal(t, 5, s.pop())
	assert.Equal(t, 4, s.pop())
	assert.Equal(t, 3, s.pop())
}
