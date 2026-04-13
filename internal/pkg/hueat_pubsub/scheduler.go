package hueat_pubsub

import (
	"github.com/hueat/backend/internal/pkg/hueat_log"
	"github.com/hueat/backend/internal/pkg/hueat_scheduler"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type pubsubScheduler struct {
	scheduler            *hueat_scheduler.Scheduler
	storage              *gorm.DB
	singleConnection     *hueat_scheduler.SingleConnection
	persistRetentionDays int
}

func newPubsubScheduler(storage *gorm.DB, scheduler *hueat_scheduler.Scheduler, persistRetentionDays int) pubsubScheduler {
	singleConnection := scheduler.GetSingleConnection(storage)
	return pubsubScheduler{
		scheduler:            scheduler,
		storage:              storage,
		singleConnection:     singleConnection,
		persistRetentionDays: persistRetentionDays,
	}
}

func (s pubsubScheduler) init() {
	// Declare all jobs to be scheduled
	var jobsToSchedule []hueat_scheduler.ScheduledJob = []hueat_scheduler.ScheduledJob{
		{
			Schedule: "0 * * * *", // Every hour at HH:00
			Handler:  s.cleanUpOldPubSubEvents,
			Parameters: hueat_scheduler.ScheduledJobParameter{
				JobID: 83701937,
				Title: "CleanUpOldPubSubEvents",
			},
		},
	}
	// Schedule all jobs
	for _, jobToSchedule := range jobsToSchedule {
		s.scheduler.AddJob(hueat_scheduler.ScheduledJob{
			Schedule:   jobToSchedule.Schedule,
			Handler:    jobToSchedule.Handler,
			Parameters: jobToSchedule.Parameters,
		})
	}

}

/*
Scheduled function to run. It cleanup expired refresh tokens
*/
func (s pubsubScheduler) cleanUpOldPubSubEvents(p hueat_scheduler.ScheduledJobParameter) error {
	defer func() {
		if r := recover(); r != nil {
			hueat_log.LogPanicError(r, "CleanUpOldPubSubEvents", "Panic occurred in cron activity")
		}
	}()
	retentionInDays := s.persistRetentionDays
	// If this istance acquires the lock, executre the business logic
	if lockAcquired := s.scheduler.AcquireLock(s.singleConnection, p.JobID); lockAcquired {
		zap.L().Info("Starting Cron Job...", zap.String("job", p.Title))
		// Delete all old events based on retention policy
		if err := s.storage.Where("event_date < NOW() - (? * INTERVAL '1 day')", retentionInDays).Delete(&eventModel{}).Error; err != nil {
			zap.L().Error("Cron Job Failed", zap.String("job", p.Title), zap.Error(err))
			return err
		}
		zap.L().Info("Cron Job executed!", zap.String("job", p.Title))
	}
	return nil
}
