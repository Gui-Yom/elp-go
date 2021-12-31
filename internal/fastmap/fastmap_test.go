package fastmap

import (
	"elp-go/internal/world"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCost(t *testing.T) {
	m := New(16, 0.9)
	m.PutCost(world.Pos(42, 69), 420.0)
	val, ok := m.GetCost(world.Pos(42, 69))
	assert.True(t, ok)
	assert.Equal(t, 420.0, val)
}

func TestPos(t *testing.T) {
	m := New(16, 0.9)
	m.PutPos(world.Pos(42, 69), world.Pos(420, 0))
	val, ok := m.GetPos(world.Pos(42, 69))
	assert.True(t, ok)
	assert.Equal(t, world.Pos(420, 0), val)
}
