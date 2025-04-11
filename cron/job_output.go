package cron

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

type SpecJob struct {
	Description string
	Spec        string
	Job         cron.Job
}

type JobOutput struct {
	fx.Out

	SpecJobs SpecJob `group:"jobs"`
}
