package datatemplate

import (
	"context"
	"strings"

	"github.com/demdxx/gocast/v2"
	"github.com/pkg/errors"
)

type IfBlock struct {
	cond      *Program
	thenBlock Block
	elseBlock Block
}

func NewIfBlock(cond *Program, thenBlock, elseBlock Block) *IfBlock {
	return &IfBlock{cond: cond, thenBlock: thenBlock, elseBlock: elseBlock}
}

func NewIfBlockWithContition(ctx context.Context, cond string, thenBlock, elseBlock Block) (Block, error) {
	cond = strings.TrimSpace(cond)
	if strings.EqualFold(cond, "true") || strings.EqualFold(cond, "1") {
		return thenBlock, nil
	}
	if strings.EqualFold(cond, "false") || strings.EqualFold(cond, "0") || strings.EqualFold(cond, "nil") || strings.EqualFold(cond, "null") {
		return elseBlock, nil
	}
	_cond, err := compileExpr(ctx, cond)
	if err != nil {
		return nil, errors.Wrap(err, cond)
	}
	return NewIfBlock(_cond, thenBlock, elseBlock), nil
}

func (b *IfBlock) String() string {
	if b.elseBlock == nil {
		return "$if: {`$expr`: `" + b.cond.Source.Content() + "`, $then: " + b.thenBlock.String() + "}"
	}
	return "$if: {`$expr`: `" + b.cond.Source.Content() + "`, $then: " + b.thenBlock.String() + ", $else: " + b.elseBlock.String() + "}"
}

func (b *IfBlock) Emit(ctx context.Context, data map[string]any) (any, error) {
	res, err := runExpr(ctx, b.cond, data)
	if err != nil {
		return nil, err
	}
	if gocast.Bool(res) {
		return b.thenBlock.Emit(ctx, data)
	}
	if b.elseBlock == nil {
		return nil, nil
	}
	return b.elseBlock.Emit(ctx, data)
}
