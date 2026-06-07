package service

import (
	"github.com/UnderTreeTech/layout/internal/dao"
	"github.com/robfig/cron/v3"
)

// Config holds the service-level configuration.
type Config struct {
	MaxIdPoolNum int
}

// Service is the core business logic layer, holding references to the DAO,
// configuration, and the cron scheduler.
type Service struct {
	cfg  *Config
	dao  dao.Dao
	cron *cron.Cron
}

// New creates and returns a new Service instance.
// It initialises the cron scheduler with panic-recovery and skip-if-still-running
// wrappers, then registers all scheduled jobs via initCrons.
func New(d dao.Dao, cfg *Config) *Service {
	// 初始化定时任务管理器
	chains := []cron.JobWrapper{cron.Recover(cron.DefaultLogger), cron.SkipIfStillRunning(cron.DefaultLogger)}
	cronManager := cron.New(cron.WithSeconds(), cron.WithChain(chains...))

	svc := &Service{
		dao:  d,
		cfg:  cfg,
		cron: cronManager,
	}

	svc.initCrons()
	return svc
}

// Close gracefully stops the cron scheduler and closes the DAO connections.
func (s *Service) Close() {
	if s.cron != nil {
		s.cron.Stop()
	}
	s.dao.Close()
}

// Ping is a no-op health-check placeholder for the service layer.
func (s *Service) Ping() {

}
