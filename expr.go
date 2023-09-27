package datatemplate

import (
	"context"
	"regexp"
	"strings"

	"github.com/demdxx/gocast/v2"
)

var reExprExtract = regexp.MustCompile(`(?mU)(?:\{\{s=|\{\{)\s*(.+)\s*\}\}`)

type ExprBlock struct {
	asStr bool
	expr  *Program
}

func NewExprBlock(expr *Program, asStr bool) *ExprBlock {
	return &ExprBlock{expr: expr, asStr: asStr}
}

func NewExprBlockFromExpr(ctx context.Context, expression string, asStr bool) (*ExprBlock, error) {
	program, err := compileExpr(ctx, expression)
	if err != nil {
		return nil, err
	}
	return NewExprBlock(program, asStr), nil
}

func (b *ExprBlock) String() string {
	return "`" + b.expr.Source.Content() + "`"
}

func (b *ExprBlock) Emit(ctx context.Context, data map[string]any) (any, error) {
	res, err := runExpr(ctx, b.expr, data)
	if err == nil && b.asStr {
		res = gocast.Str(res)
	}
	return res, err
}

type ExprBlockStringTmplate struct {
	expression string
	exprs      map[string]*Program
}

func NewExprBlockFromString(ctx context.Context, data string) (any, error) {
	matches := reExprExtract.FindAllStringSubmatch(data, -1)
	if len(matches) == 0 {
		return data, nil
	}
	if len(matches) == 1 && matches[0][0] == data {
		return NewExprBlockFromExpr(ctx, matches[0][1], strings.HasPrefix(matches[0][0], "{{s="))
	}
	exprs := make(map[string]*Program, len(matches))
	for _, match := range matches {
		if exprs[match[0]] != nil {
			continue
		}
		program, err := compileExpr(ctx, match[1])
		if err != nil {
			return nil, err
		}
		exprs[match[0]] = program
	}
	return &ExprBlockStringTmplate{expression: data, exprs: exprs}, nil
}

func (b *ExprBlockStringTmplate) String() string {
	return b.expression
}

func (b *ExprBlockStringTmplate) Emit(ctx context.Context, data map[string]any) (any, error) {
	result := b.expression
	for k, v := range b.exprs {
		res, err := runExpr(ctx, v, data)
		if err != nil {
			return nil, err
		}
		result = strings.ReplaceAll(result, k, gocast.Str(res))
	}
	return result, nil
}
