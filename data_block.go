package datatemplate

import (
	"bytes"
	"context"
	"fmt"

	"github.com/demdxx/gocast/v2"
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

func (b *DataBlock) String() string {
	if sp, _ := b.data.(fmt.Stringer); sp != nil {
		return sp.String()
	}
	return gocast.Str(b.data)
}

func (b *DataBlock) Emit(ctx context.Context, data map[string]any) (any, error) {
	return b.data, nil
}

type DataBlockSlice struct {
	data []any
}

func (b *DataBlockSlice) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, item := range b.data {
		if i > 0 {
			_, _ = buf.WriteString(", ")
		}
		if sp, _ := item.(fmt.Stringer); sp != nil {
			_, _ = buf.WriteString(sp.String())
		} else {
			_, _ = buf.WriteString(gocast.Str(item))
		}
	}
	buf.WriteString("]")
	return buf.String()
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

func (b *DataBlockMap) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	i := 0
	for key, item := range b.data {
		if i > 0 {
			_, _ = buf.WriteString(", ")
		}
		_, _ = buf.WriteString(key)
		_, _ = buf.WriteString(": ")
		if sp, _ := item.(fmt.Stringer); sp != nil {
			_, _ = buf.WriteString(sp.String())
		} else {
			_, _ = buf.WriteString(gocast.Str(item))
		}
		i++
	}
	buf.WriteString("}")
	return buf.String()
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
