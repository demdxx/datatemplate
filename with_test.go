package datatemplate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithBlock(t *testing.T) {
	exprStr, err := NewExprBlockFromExpr(context.Background(), "np.name", true)
	assert.NoError(t, err)

	withBlock, err := NewWithBlockFromExpr(context.TODO(), "np", "person[0]", exprStr)
	assert.NoError(t, err)

	res, err := withBlock.Emit(context.TODO(), map[string]any{
		"person": []map[string]any{
			{
				"name": "tony",
				"age":  42,
			},
			{
				"name": "peter",
				"age":  16,
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "tony", res)
}
