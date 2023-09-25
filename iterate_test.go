package datatemplate

import (
	"context"
	"testing"

	"github.com/demdxx/gocast/v2"
	"github.com/stretchr/testify/assert"
)

func TestIterateBlock(t *testing.T) {
	ctx := context.Background()
	data := map[string]any{
		"data": map[string]any{
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
		},
	}
	tests := []struct {
		iterator string
		block    Block
		res      []any
	}{
		{
			iterator: "data.person",
			block:    tNewExpr(ctx, "{{item.name}}"),
			res:      []any{"tony", "peter"},
		},
		{
			iterator: "data.person",
			block:    tNewExpr(ctx, "{{item.age}}"),
			res:      []any{42, 16},
		},
		{
			iterator: "data.person",
			block:    tNewExpr(ctx, "{{item.name}} is {{item.age}} years old"),
			res:      []any{"tony is 42 years old", "peter is 16 years old"},
		},
		{
			iterator: "data.person[0]",
			block:    tNewExpr(ctx, "{{key}}: {{item}}"),
			res:      []any{"name: tony", "age: 42"},
		},
	}

	for i, test := range tests {
		t.Run(gocast.Str(i), func(t *testing.T) {
			iterate, err := NewIterateBlockFromExpr(ctx, test.iterator, "", "", "", test.block)
			if err != nil {
				t.Fatal(err)
			}

			res, err := iterate.Emit(context.TODO(), data)
			if err != nil {
				t.Fatal(err)
			}
			assert.ElementsMatch(t, test.res, res)
		})
	}
}

func tNewExpr(ctx context.Context, tpl string) Block {
	bl, _ := NewExprBlockFromString(ctx, tpl)
	return bl.(Block)
}
