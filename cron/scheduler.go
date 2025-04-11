package cron

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"

	"github.com/enesanbar/go-service/log"
	"github.com/enesanbar/go-service/wiring"
)

type Scheduler struct {
	cron     *cron.Cron
	logger   log.Factory
	specJobs []SpecJob
}

type SchedulerParams struct {
	fx.In

	SpecJobs []SpecJob `group:"jobs"`
	Logger   log.Factory
}

func NewScheduler(p SchedulerParams) (wiring.RunnableGroup, *Scheduler) {
	scheduler := &Scheduler{cron: cron.New(), logger: p.Logger, specJobs: p.SpecJobs}
	return wiring.RunnableGroup{Runnable: scheduler}, scheduler
}

func (s *Scheduler) Start() error {
	s.logger.Bg().Info("Getting all registered CRON jobs...")
	for _, job := range s.specJobs {
		s.logger.Bg().Infof("[%s] Registering job in the scheduler", job.Description)
		entryID, err := s.cron.AddJob(job.Spec, job.Job)
		if err != nil {
			s.logger.Bg().Infof("[%s] Unable to register the job", job.Description)
			continue
		}
		s.logger.Bg().Infof("[%s] Job has been registered in the scheduler with id '%d'.", job.Description, entryID)
	}

	s.logger.Bg().Info("All jobs has been registered. Starting cron scheduler...")
	s.cron.Start()
	// s.logger.Infof("Running jobs: %+v", s.cron.Entries())
	return nil
}

func (s *Scheduler) Stop() error {
	s.logger.Bg().Info("Stopping cron scheduler...")
	s.cron.Stop()
	return nil
}
