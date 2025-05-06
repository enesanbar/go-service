package mysql

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"persistence.mysql",
	fx.Provide(MySQLConnections),
)
