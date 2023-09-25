package datatemplate

import (
	"context"
	"testing"

	"github.com/demdxx/gocast/v2"
	"github.com/stretchr/testify/assert"
)

func TestExpr(t *testing.T) {
	ctx := context.Background()

	t.Run("expr", func(t *testing.T) {
		data := map[string]any{"a": 1}
		tests := []struct {
			exp     string
			res     any
			expType any
			fnk     func(context.Context, string) (any, error)
		}{
			{"1", 1, &ExprBlock{}, func(ctx context.Context, exp string) (any, error) { return NewExprBlockFromExpr(ctx, exp, false) }},
			{"1 + 1", 2, &ExprBlock{}, func(ctx context.Context, exp string) (any, error) { return NewExprBlockFromExpr(ctx, exp, false) }},
			{"a + 1", 2, &ExprBlock{}, func(ctx context.Context, exp string) (any, error) { return NewExprBlockFromExpr(ctx, exp, false) }},
			{"a * 3 - 1", 2, &ExprBlock{}, func(ctx context.Context, exp string) (any, error) { return NewExprBlockFromExpr(ctx, exp, false) }},
			{"{{1}}", 1, &ExprBlock{}, NewExprBlockFromString},
			{"{{s=1}}", "1", &ExprBlock{}, NewExprBlockFromString},
			{"{{1 + 1}}", 2, &ExprBlock{}, NewExprBlockFromString},
			{"{{a + 1}}", 2, &ExprBlock{}, NewExprBlockFromString},
		}

		for i, test := range tests {
			t.Run(gocast.Str(i), func(t *testing.T) {
				exp, err := test.fnk(ctx, test.exp)
				if !assert.NoError(t, err) || !assert.IsType(t, test.expType, exp) {
					return
				}
				res, err := exp.(Block).Emit(context.Background(), data)
				assert.NoError(t, err)
				assert.Equal(t, test.res, res)
			})
		}
	})

	t.Run("exprStemplate", func(t *testing.T) {
		data := map[string]any{"a": 1}
		tests := []struct {
			exp     string
			res     any
			expType any
			fnk     func(context.Context, string) (any, error)
		}{
			{"{{s=1}}", "1", &ExprBlock{}, NewExprBlockFromString},
			{"{{s=100*200}}", "20000", &ExprBlock{}, NewExprBlockFromString},
			{"X: {{1}}", "X: 1", &ExprBlockStringTmplate{}, NewExprBlockFromString},
			{"X: {{1 + 1}}", "X: 2", &ExprBlockStringTmplate{}, NewExprBlockFromString},
			{"X: {{a + 1}}", "X: 2", &ExprBlockStringTmplate{}, NewExprBlockFromString},
			{"{{a + 1}} {{a + 1}}", "2 2", &ExprBlockStringTmplate{}, NewExprBlockFromString},
			{"{{a + 1}} {{a + 1}} {{a + 1}}", "2 2 2", &ExprBlockStringTmplate{}, NewExprBlockFromString},
			{"{{a + 1}} {{a + a}} {{a * 2}} {{a*4 / 2}}", "2 2 2 2", &ExprBlockStringTmplate{}, NewExprBlockFromString},
		}

		for i, test := range tests {
			t.Run(gocast.Str(i), func(t *testing.T) {
				exp, err := test.fnk(ctx, test.exp)
				if !assert.NoError(t, err) || !assert.IsType(t, test.expType, exp) {
					return
				}
				res, err := exp.(Block).Emit(context.Background(), data)
				assert.NoError(t, err)
				assert.Equal(t, test.res, res)
			})
		}
	})
}
