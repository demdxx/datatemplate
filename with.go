package datatemplate

import (
	"context"

	"github.com/demdxx/xtypes"
)

type WithBlock struct {
	name string
	expr *Program
	body Block
}

func NewWithBlock(name string, expr *Program, body Block) *WithBlock {
	return &WithBlock{name: name, expr: expr, body: body}
}

func NewWithBlockFromExpr(ctx context.Context, name, expr string, body Block) (Block, error) {
	_expr, err := compileExpr(ctx, expr)
	if err != nil {
		return nil, err
	}
	return NewWithBlock(name, _expr, body), nil
}

func (wi *WithBlock) String() string {
	return "$with: {`$expr`: `" + wi.name + " := " + wi.expr.Source.Content() + "`, $body: " + wi.body.String() + "}"
}

func (wi *WithBlock) Emit(ctx context.Context, data map[string]any) (any, error) {
	res, err := runExpr(ctx, wi.expr, data)
	if err != nil {
		return nil, err
	}
	newData := xtypes.Map[string, any](data).Copy().Set(wi.name, res)
	return wi.body.Emit(ctx, newData)
}
