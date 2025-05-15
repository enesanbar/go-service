package middlewares

import "go.uber.org/fx"

func AsMiddleware(p any) any {
	return fx.Annotate(
		p,
		fx.ResultTags(`group:"middlewares"`),
	)
}
