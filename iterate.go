package datatemplate

import (
	"context"

	"github.com/demdxx/gocast/v2"
	"github.com/demdxx/xtypes"
	"github.com/pkg/errors"
)

var errInvalidIteratator = errors.New("invalid iterator")

type IterateBlock struct {
	expr      *Program
	indexName string
	keyName   string
	valueName string
	block     Block
}

func NewIterateBlock(expr *Program, indexName, keyName, valueName string, block Block) *IterateBlock {
	return &IterateBlock{
		expr:      expr,
		indexName: strOrDef(indexName, "index"),
		keyName:   strOrDef(keyName, "key"),
		valueName: strOrDef(valueName, "item"),
		block:     block,
	}
}

func NewIterateBlockFromExpr(ctx context.Context, expression, indexName, keyName, valueName string, block Block) (*IterateBlock, error) {
	program, err := compileExpr(ctx, expression)
	if err != nil {
		return nil, errors.Wrap(err, expression)
	}
	return NewIterateBlock(program, indexName, keyName, valueName, block), nil
}

func (it *IterateBlock) String() string {
	return "$iterate: {`$expr`: `" + it.expr.Source.Content() +
		"`, $index: `" + it.indexName +
		"`, $key: `" + it.keyName +
		"`, $value: `" + it.valueName +
		"`, $body: " + it.block.String() + "}"
}

func (it *IterateBlock) Emit(ctx context.Context, data map[string]any) (any, error) {
	otData, err := runExpr(ctx, it.expr, data)
	if err != nil {
		return nil, err
	}
	if !gocast.IsSlice(otData) && !gocast.IsMap(otData) {
		return nil, errors.Wrap(errInvalidIteratator, "not a slice or map")
	}

	// copy context data
	nData := xtypes.Map[string, any](data).Copy()

	// iterate slice object data
	if gocast.IsSlice(otData) {
		list := gocast.AnySlice[any](otData)
		res := make([]any, 0, len(list))
		for index, item := range list {
			nData[it.indexName] = index
			nData[it.valueName] = item
			if rData, err := it.block.Emit(ctx, nData); err != nil {
				return nil, err
			} else if rData != nil {
				res = append(res, rData)
			}
		}
		return res, nil
	}

	// iterate map object data
	mp := gocast.Map[string, any](otData)
	res := make([]any, 0, len(mp))
	index := 0
	for key, item := range mp {
		nData[it.keyName] = key
		nData[it.valueName] = item
		nData[it.indexName] = index
		if rData, err := it.block.Emit(ctx, nData); err != nil {
			return nil, err
		} else {
			res = append(res, rData)
		}
		index++
	}
	return res, nil
}

func strOrDef(s, def string) string {
	if s != "" {
		return s
	}
	return def
}
