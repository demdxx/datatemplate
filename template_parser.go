package datatemplate

import (
	"context"
	"regexp"
	"strings"

	"github.com/demdxx/gocast/v2"
	"github.com/demdxx/xtypes"
	"github.com/pkg/errors"
)

var (
	errInvalidIfBlock                        = errors.New("invalid if block")
	errInvalidIteratorBlock                  = errors.New("invalid iterator block")
	errInvalidWithBlock                      = errors.New("invalid with block")
	errInvalidWithBlockExpr                  = errors.New("invalid with block expr")
	errDataFieldsIsNotAllowedIfBodyIsDefined = errors.New("data fields is not allowed if body is defined")
)

func parseBlocks(ctx context.Context, data any) (any, error) {
	switch {
	case gocast.IsSlice(data):
		arr := gocast.AnySlice[any](data)
		blocks := make([]any, 0, len(arr))
		hasBlocks := false
		for _, item := range arr {
			block, err := parseBlocks(ctx, item)
			if err != nil {
				return nil, err
			}
			if _, ok := block.(Block); ok {
				hasBlocks = true
			}
			blocks = append(blocks, block)
		}
		if hasBlocks {
			return &DataBlockSlice{data: blocks}, nil
		}
	case gocast.IsMap(data) || gocast.IsStruct(data):
		m := gocast.Map[string, any](data)

		if _, ok := m["$if"]; ok {
			return parseIfBlock(ctx, m)
		}
		if _, ok := m["$iterate"]; ok {
			return parseIteratorBlock(ctx, m)
		}
		if _, ok := m["$with"]; ok {
			return parseWithBlock(ctx, m)
		}

		blocks := make(map[string]any, len(m))
		hasBlocks := false
		for key, item := range m {
			block, err := parseBlocks(ctx, item)
			if err != nil {
				return nil, err
			}
			if _, ok := block.(Block); ok {
				hasBlocks = true
			}
			blocks[key] = block
		}
		if hasBlocks {
			return &DataBlockMap{data: blocks}, nil
		}
	case gocast.IsStr(data):
		return NewExprBlockFromString(ctx, gocast.Str(data))
	}
	return data, nil
}

// Example 1:
// $if: "person.age > 18"
// field1: "value1"
// field2: "value2"
//
// Example 2:
// $if:
//
//	$cond: "person.age > 18"
//	field1: "value1"
//	field2: "value2"
//
// $else:
//
//	field1: "value1"
//	field2: "value2"
func parseIfBlock(ctx context.Context, data map[string]any) (Block, error) {
	var (
		ifdata, ok = data["$if"]
		condition  string
		thenBlock  Block
		elseBlock  Block
	)
	if !ok {
		return nil, errInvalidIfBlock
	}

	if gocast.IsStr(ifdata) {
		condition = gocast.Str(ifdata)
		dataCopy := xtypes.Map[string, any](data).Filter(func(k string, _ any) bool { return k != "$if" })
		body, err := parseBlocks(ctx, dataCopy)
		if err != nil {
			return nil, err
		}
		thenBlock = NewDataBlock(body)
	} else {
		condData := xtypes.Map[string, any](gocast.Map[string, any](ifdata)).Copy()
		condition = gocast.Str(condData["$cond"])
		if condition == "" {
			condition = gocast.Str(condData["$condition"])
		}
		delete(condData, "$cond")
		delete(condData, "$condition")
		body, err := parseBlocks(ctx, condData)
		if err != nil {
			return nil, err
		}
		thenBlock = NewDataBlock(body)
		body, err = parseBlocks(ctx, data["$else"])
		if err != nil {
			return nil, err
		}
		elseBlock = NewDataBlock(body)
	}

	return NewIfBlockWithContition(ctx, condition, thenBlock, elseBlock)
}

// Example 1:
// $iterate: "data.list"
// field1: "{{item.field1}}"
// field2: "{{item.field2}}"
// index:  "{{index}}"
//
// Example 2:
// $iterate:
//
//	$expr: "data.list"
//	$body: "{{index}}"
//
// Example 3:
// $iterate:
//
//	$expr: "data.list"
//	field1: "{{item.field1}}"
//	field2: "{{item.field2}}"
//	index:  "{{index}}"
//
// Example 4:
// $iterate: "data.list"
// $body: "{{index}}"
func parseIteratorBlock(ctx context.Context, data map[string]any) (Block, error) {
	var (
		iterateData, ok = data["$iterate"]
		iterateExpr     string
		bodyData        any
	)
	if !ok {
		return nil, errInvalidIteratorBlock
	}

	if gocast.IsStr(iterateData) {
		iterateExpr = gocast.Str(iterateData)

		// If body is defined then we should not have any other fields
		if bodyData, ok = data["$body"]; ok {
			if len(data) > 2 {
				return nil, errDataFieldsIsNotAllowedIfBodyIsDefined
			}
		} else {
			// Remove $iterate field if present
			bodyData = xtypes.Map[string, any](data).Filter(func(k string, _ any) bool { return k != "$iterate" })
		}
	} else {
		dataCopy := xtypes.Map[string, any](gocast.Map[string, any](iterateData)).Copy()
		iterateExpr = gocast.Str(dataCopy["$expr"])

		// If body is defined then we should not have any other fields
		if bodyData, ok = dataCopy["$body"]; ok {
			if len(dataCopy) > 2 {
				return nil, errDataFieldsIsNotAllowedIfBodyIsDefined
			}
		} else {
			delete(dataCopy, "$expr")
			bodyData = dataCopy
		}
	}

	// Parse blocks from data
	body, err := parseBlocks(ctx, bodyData)
	if err != nil {
		return nil, err
	}

	return NewIterateBlockFromExpr(ctx, iterateExpr, "", "", "", NewDataBlock(body))
}

// Extract variable name from expression like: varName := expr
var reLeftVariableName = regexp.MustCompile(`^\s*([a-zA-Z0-9_]+)\s*:=\s*`)

// Example 1:
// $with: varName := np.name
// $body:
//
//	field1: "{{var}}"
//	field2: "{{np.age}}"
//
// Example 2:
// $with:
//
//	$expr: varName := np.name
//	$body:
//		field1: "{{var}}"
//		field2: "{{np.age}}"
//
// Example 3:
// $with:
//
//	$expr: varName := np.name
//	field1: "{{var}}"
//	field2: "{{np.age}}"
//
// Example 4:
// $with: varName := np.name
// field1: "{{var}}"
// field2: "{{np.age}}"
func parseWithBlock(ctx context.Context, data map[string]any) (Block, error) {
	var (
		withData, ok = data["$with"]
		withExpr     string
		bodyData     any
	)
	if !ok {
		return nil, errInvalidWithBlock
	}

	if gocast.IsStr(withData) {
		withExpr = gocast.Str(withData)

		// If body is defined then we should not have any other fields
		if bodyData, ok = data["$body"]; ok {
			if len(data) > 2 {
				return nil, errDataFieldsIsNotAllowedIfBodyIsDefined
			}
		} else {
			// Remove $with field if present
			bodyData = xtypes.Map[string, any](data).Filter(func(k string, _ any) bool { return k != "$with" })
		}
	} else {
		dataCopy := xtypes.Map[string, any](gocast.Map[string, any](withData)).Copy()
		withExpr = gocast.Str(dataCopy["$expr"])

		// If body is defined then we should not have any other fields
		if bodyData, ok = dataCopy["$body"]; ok {
			if len(dataCopy) > 2 {
				return nil, errDataFieldsIsNotAllowedIfBodyIsDefined
			}
		} else {
			delete(dataCopy, "$expr")
			bodyData = dataCopy
		}
	}

	varArr := reLeftVariableName.FindStringSubmatch(withExpr)
	if len(varArr) < 2 {
		return nil, errInvalidWithBlockExpr
	}
	withExpr = strings.TrimSpace(strings.Replace(withExpr, varArr[0], "", 1))

	// Parse blocks from data
	body, err := parseBlocks(ctx, bodyData)
	if err != nil {
		return nil, err
	}

	return NewWithBlockFromExpr(ctx, varArr[1], withExpr, NewDataBlock(body))
}
