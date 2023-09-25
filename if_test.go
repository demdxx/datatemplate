package datatemplate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfBlock(t *testing.T) {
	ctx := context.Background()

	t.Run("if", func(t *testing.T) {
		_if, err := NewIfBlockWithContition(ctx, "person.name == 'tony'", NewDataBlock("true"), NewDataBlock("false"))
		assert.NoError(t, err)

		res, err := _if.Emit(context.TODO(), map[string]any{"person": map[string]any{"name": "tony"}})
		assert.NoError(t, err)
		assert.Equal(t, "true", res)

		res, err = _if.Emit(context.TODO(), map[string]any{"person": map[string]any{"name": "peter"}})
		assert.NoError(t, err)
		assert.Equal(t, "false", res)
	})
}
