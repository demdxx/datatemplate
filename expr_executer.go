package datatemplate

import (
	"context"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

var ctxExprKey = struct{}{}

func ctxWithExprOptions(ctx context.Context, options ...expr.Option) context.Context {
	return context.WithValue(ctx, ctxExprKey, options)
}

func ctxExprOptions(ctx context.Context) []expr.Option {
	if options, ok := ctx.Value(ctxExprKey).([]expr.Option); ok {
		return options
	}
	return nil
}

type Program = vm.Program

func compileExpr(ctx context.Context, expression string) (*Program, error) {
	return expr.Compile(expression, ctxExprOptions(ctx)...)
}
