// Package: datatemplate represents data template engine which can be used
// to generate any data constructions from data template.
package datatemplate

import (
	"context"
	"fmt"
)

type Block interface {
	fmt.Stringer
	Emit(ctx context.Context, data map[string]any) (any, error)
}

type Template struct {
	root Block
}

// NewTemplate creates new template from root block
func NewTemplate(root Block) *Template {
	return &Template{root: root}
}

// NewTemplateFor creates new template from data input (string, map, struct, etc)
func NewTemplateFor(data any, opts ...Option) (*Template, error) {
	var opt options
	for _, o := range opts {
		o(&opt)
	}
	root, err := parseBlocks(ctxWithExprOptions(context.Background(), opt.exprOpts...), data)
	if err != nil {
		return nil, err
	}
	return NewTemplate(NewDataBlock(root)), nil
}

func (tpl *Template) String() string {
	return tpl.root.String()
}

// Process template with data and return result according to template of data
func (tpl *Template) Process(ctx context.Context, data map[string]any) (any, error) {
	return tpl.root.Emit(ctx, data)
}
