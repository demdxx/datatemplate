package datatemplate

import "github.com/antonmedv/expr"

type options struct {
	exprOpts []expr.Option
}

type Option func(o *options)

// WithExprOptions sets expr options for the module "github.com/antonmedv/expr"
func WithExprOptions(opts ...expr.Option) Option {
	return func(o *options) {
		o.exprOpts = append(o.exprOpts, opts...)
	}
}

// WithExprEnv sets environment for expressions
func WithExprEnv(env any) Option {
	return WithExprOptions(expr.Env(env))
}
