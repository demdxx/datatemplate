package datatemplate

import (
	"context"
)

type DataBlock struct {
	data any
}

func NewDataBlock(data any) Block {
	switch b := data.(type) {
	case nil:
		return nil
	case Block:
		return b
	default:
		return &DataBlock{data: data}
	}
}

func (b *DataBlock) Emit(ctx context.Context, data map[string]any) (any, error) {
	return b.data, nil
}

type DataBlockSlice struct {
	data []any
}

func (b *DataBlockSlice) Emit(ctx context.Context, data map[string]any) (any, error) {
	newResult := make([]any, 0, len(b.data))
	for _, item := range b.data {
		switch b := item.(type) {
		case Block:
			res, err := b.Emit(ctx, data)
			if err != nil {
				return nil, err
			}
			newResult = append(newResult, res)
		default:
			newResult = append(newResult, item)
		}
	}
	return newResult, nil
}

type DataBlockMap struct {
	data map[string]any
}

func (b *DataBlockMap) Emit(ctx context.Context, data map[string]any) (any, error) {
	newResult := make(map[string]any, len(b.data))
	for key, item := range b.data {
		switch b := item.(type) {
		case Block:
			res, err := b.Emit(ctx, data)
			if err != nil {
				return nil, err
			}
			newResult[key] = res
		default:
			newResult[key] = item
		}
	}
	return newResult, nil
}
