package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_NewFromString(t *testing.T) {
	var identifier = "9BO"
	var nonzero uint64 = 39282

	b := Builder{}
	c, err := b.NewFromString(identifier)

	require.NoError(t, err)
	assert.Equal(t, c.ID(), nonzero)
	assert.Equal(t, c.String(), identifier)
}

func TestBuilder_NewFromID(t *testing.T) {
	var identifier = "9BO"
	var nonzero uint64 = 39282

	b := Builder{}
	c, err := b.NewFromID(nonzero)

	require.NoError(t, err)
	assert.Equal(t, c.ID(), nonzero)
	assert.Equal(t, c.String(), identifier)
}

func TestBuilder_New(t *testing.T) {
	var zero uint64

	b := Builder{}
	c0 := b.New()
	c1 := b.New()

	assert.Greater(t, c0.ID(), zero)
	assert.Greater(t, c1.ID(), zero)
	assert.NotEqual(t, c0.ID(), c1.ID())
}
