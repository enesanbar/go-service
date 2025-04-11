package service

import (
	"github.com/enesanbar/go-service/validation"
	"go.uber.org/fx"
)

// Module combines other modules to bootstrap the server
var Module = fx.Options(
	// rest.Module,
	// router.Module,
	validation.Module,
)
