package datatemplate

import (
	"context"
	"reflect"
	"testing"

	"github.com/demdxx/gocast/v2"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	commonData := map[string]any{
		"name":    "tony",
		"surname": "stark",
		"age":     42,
		"person": []map[string]any{
			{
				"name": "tony",
				"age":  42,
			},
			{
				"name": "rony",
				"age":  14,
			},
		},
	}

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			tpl string
			res string
		}{
			{"Hello {{name}}", "Hello tony"},
			{"Hello {{name}} {{surname}}", "Hello tony stark"},
			{"Hello {{name}} {{surname}} {{age}}", "Hello tony stark 42"},
			{"Hello {{name}} {{surname}} {{age}} {{person[0].name}}", "Hello tony stark 42 tony"},
			{"Hello {{name}} {{surname}} {{age}} {{person[0].name}} {{person[0].age}}", "Hello tony stark 42 tony 42"},
			{"Hello {{name}} {{surname}} {{age}} {{person[0].name}} {{person[0].age}} {{person[0].age > 18 ? 'adult' : 'teenager'}}", "Hello tony stark 42 tony 42 adult"},
		}
		for i, test := range tests {
			t.Run(gocast.Str(i), func(t *testing.T) {
				tpl, err := NewTemplateFor(test.tpl)
				if !assert.NoError(t, err) {
					return
				}
				res, err := tpl.Process(context.TODO(), commonData)
				assert.NoError(t, err)
				assert.Equal(t, test.res, res)
			})
		}
	})

	t.Run("map", func(t *testing.T) {
		tests := []struct {
			tpl map[string]any
			res map[string]any
		}{
			{
				tpl: map[string]any{
					"person": map[string]any{
						"name": "{{person[0].name}}",
						"age":  "{{person[0].age}}",
					},
				},
				res: map[string]any{
					"person": map[string]any{
						"name": "tony",
						"age":  42,
					},
				},
			},
			{
				tpl: map[string]any{
					"person": "{{person[0]}}",
					"desc":   "Hello {{person[0].name}} of {{person[0].age}} years old",
					"what":   "Done!",
				},
				res: map[string]any{
					"person": map[string]any{
						"name": "tony",
						"age":  42,
					},
					"desc": "Hello tony of 42 years old",
					"what": "Done!",
				},
			},
			{
				tpl: map[string]any{
					"list_fields": []any{"Name: {{person[0].name}}", "Age: {{person[0].age}}", "Good!"},
				},
				res: map[string]any{
					"list_fields": []any{"Name: tony", "Age: 42", "Good!"},
				},
			},
			// If statement tests
			{
				tpl: map[string]any{
					"person": map[string]any{
						"$if":  "person[0].age > 18",
						"name": "{{person[0].name}}",
						"age":  "{{person[0].age}}",
					},
				},
				res: map[string]any{
					"person": map[string]any{
						"name": "tony",
						"age":  42,
					},
				},
			},
			{
				tpl: map[string]any{
					"person": map[string]any{
						"$if":  "person[0].age < 18",
						"name": "{{person[0].name}}",
						"age":  "{{person[0].age}}",
					},
				},
				res: map[string]any{
					"person": nil,
				},
			},
			{
				tpl: map[string]any{
					"person": map[string]any{
						"$if": map[string]any{
							"$cond": "person[0].age < 18",
							"name":  "{{person[0].name}}",
							"age":   "{{person[0].age}}",
						},
						"$else": map[string]any{
							"name": "{{person[1].name}}",
							"age":  "{{person[1].age}}",
						},
					},
				},
				res: map[string]any{
					"person": map[string]any{
						"name": "rony",
						"age":  14,
					},
				},
			},
			// Iterate statement tests
			{
				tpl: map[string]any{
					"persons": map[string]any{
						"$iterate": "person",
						"index":    "{{index}}",
						"name":     "{{item.name}}",
						"age":      "{{item.age}}",
					},
				},
				res: map[string]any{
					"persons": []any{
						map[string]any{
							"index": 0,
							"name":  "tony",
							"age":   42,
						},
						map[string]any{
							"index": 1,
							"name":  "rony",
							"age":   14,
						},
					},
				},
			},
			{
				tpl: map[string]any{
					"names": map[string]any{
						"$iterate": map[string]any{
							"$expr": "person",
							"$body": "{{item.name}}",
						},
					},
				},
				res: map[string]any{
					"names": []any{"tony", "rony"},
				},
			},
			{
				tpl: map[string]any{
					"names": map[string]any{
						"$iterate": "person",
						"$body":    "{{item.name}}",
					},
				},
				res: map[string]any{
					"names": []any{"tony", "rony"},
				},
			},
			{
				tpl: map[string]any{
					"persons": map[string]any{
						"$iterate": "person",
						"$body": map[string]any{
							"index": "{{s= index + 1}}",
							"name":  "{{item.name}}",
							"age":   "{{s= item.age}}",
						},
					},
				},
				res: map[string]any{
					"persons": []any{
						map[string]any{
							"index": "1",
							"name":  "tony",
							"age":   "42",
						},
						map[string]any{
							"index": "2",
							"name":  "rony",
							"age":   "14",
						},
					},
				},
			},
		}
		for i, test := range tests {
			t.Run(gocast.Str(i), func(t *testing.T) {
				tpl, err := NewTemplateFor(test.tpl)
				if !assert.NoError(t, err) {
					return
				}
				res, err := tpl.Process(context.TODO(), commonData)
				assert.NoError(t, err)
				if !assert.True(t, reflect.DeepEqual(test.res, res)) {
					t.Logf("Expected: %#v", test.res)
					t.Logf("Result: %#v", res)
				}
			})
		}
	})
}
