// Package jobs provides cron job-related functionalities.
package jobs

import (
	"log"

	"gin-alpine/src/internal/configs"
	"gin-alpine/src/pkg/utils"

	"go.uber.org/zap"

	"github.com/robfig/cron/v3"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Job struct {
	Cron          *cron.Cron
	Logger        *zap.Logger
	Config        *configs.Config
	RollingLogger *lumberjack.Logger
}

func NewJob(
	cron *cron.Cron,
	logger *zap.Logger,
) *Job {
	return &Job{
		Cron:   cron,
		Logger: logger,
	}
}

func (j *Job) RotateLogs() {
	_, err := j.Cron.AddFunc("@midnight", func() {
		j.Logger.Info(utils.JobStarted, zap.String("job", utils.JobRotateLogs))
		err := j.RollingLogger.Rotate()
		if err != nil {
			j.Logger.Error(
				"job error",
				zap.String("job", utils.JobRotateLogs),
				zap.String("message", err.Error()),
			)
			j.Logger.Info(utils.JobFinished, zap.String("job", utils.JobRotateLogs))
		}
		j.Logger.Info(utils.JobFinished, zap.String("job", utils.JobRotateLogs))
	})
	if err != nil {
		log.Printf("error scheduling the job %v\n", err)
		return
	}
}

func (j *Job) RegisterError(err error, desc string) {
	j.Logger.Error(
		"job error",
		zap.String("job", desc),
		zap.String("message", err.Error()),
	)
}
