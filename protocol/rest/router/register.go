package router

import "go.uber.org/fx"

func AsRoute(p any) any {
	return fx.Annotate(
		p,
		fx.ResultTags(`group:"routes"`),
	)
}
